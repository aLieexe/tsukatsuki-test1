package main

import (
	"go-webserver/internal/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	// new instance off app, this will just discard anything written to it
	app := newTestApplication(t)

	// create a new test server, passing in the app.routes()
	// DONT FORGET routes is a function not a field idiot
	// thats why we are able to access it
	// also this server will run in random port, we dont need to know
	server := newTestServer(t, app.routes())

	//make sure to close the server
	defer server.Close()

	//make a new get request to url/ping,
	statusCode, _, body := server.get(t, "/ping")

	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, string(body), "OKE")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/snippet-123",
			wantCode: http.StatusOK,
			wantBody: "RIO RIO RIO RIO RIO ",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/snippet-1234121",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := server.get(t, tt.urlPath)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()
	_, _, body := server.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)
	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid email",
			userName:     validName,
			userEmail:    "bob@example.",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pa$$",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate email",
			userName:     validName,
			userEmail:    "dupe@example.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)
			code, _, body := server.postForm(t, "/user/signup", form)
			assert.Equal(t, code, tt.wantCode)
			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	server := newTestServer(t, app.routes())
	defer server.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		statusCode, headers, _ := server.get(t, "/snippet/create")
		assert.Equal(t, statusCode, 303)

		//we do to the /user/login, because we get redirected, which is why not /snippet/create
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		//yoink the csrf
		_, _, body := server.get(t, "/user/login")
		validCSRFToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", validCSRFToken)

		server.postForm(t, "/user/login", form)

		statusCode, _, body := server.get(t, "/snippet/create")

		assert.Equal(t, statusCode, http.StatusOK)
		assert.StringContains(t, body, `<form action="/snippet/create" method="POST">`)
	})
}

// func TestSnippetCreate(t *testing.T) {
// 	app := newTestApplication(t)
// 	server := newTestServer(t, app.routes())

// 	defer server.Close()
// 	_, _, body := server.get(t, "/snippet/create")
// 	validCSRFToken := extractCSRFToken(t, body)

// 	const (
// 		validTitle      = "Tsukatsuki Rio"
// 		validContent    = "RIO RIO RIO"
// 		validExpiration = 1
// 		formTag         = "<form action='/snippet/create' method='POST' novalidate>"
// 	)

// 	tests := []struct {
// 		name       string
// 		title      string
// 		content    string
// 		expiration int
// 		csrfToken  string
// 		wantCode   int
// 	}{
// 		{
// 			name:       "valid",
// 			title:      validTitle,
// 			content:    validContent,
// 			expiration: validExpiration,
// 			csrfToken:  validCSRFToken,
// 			wantCode:   200,
// 		},
// 	}

// 	for _, test := range tests {
// 		fmt.Println("KLJSDFKLJDFSKLJSDKLJfkj")
// 		t.Run(test.name, func(t *testing.T) {
// 			form := url.Values{}
// 			form.Add("title", test.title)
// 			form.Add("content", test.content)
// 			form.Add("expiration", "2")
// 			form.Add("csrf_token", test.csrfToken)
// 			code, _, _ := server.postForm(t, "/snippet/create", form)
// 			assert.Equal(t, code, test.wantCode)

// 		})
// 	}

// }
