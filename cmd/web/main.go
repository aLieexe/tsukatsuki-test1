package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"go-webserver/internal/models"
	"go-webserver/utils"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type application struct {
	debug          bool
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func connectDb() (*pgxpool.Pool, error) {
	stringConnection := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v",
		utils.GetEnv("PGUSER"), utils.GetEnv("PGPASSWORD"), utils.GetEnv("PGHOST"), utils.GetEnv("PGPORT"), utils.GetEnv("PGDATABASE"))
	config, err := pgxpool.ParseConfig(stringConnection)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	if err := godotenv.Load(".env"); err != nil {
		logger.Error("Error loading .env file")
	}

	db, err := connectDb()
	if err != nil {
		logger.Error(err.Error())
	}

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = time.Hour * 12
	sessionManager.Cookie.Secure = true

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
	}

	debugEnv := utils.GetEnv("DEBUG", "false")
	var debug bool
	if debugEnv == "true" {
		debug = true
	} else {
		debug = false
	}

	fmt.Println(debug)

	app := &application{
		debug:          debug,
		logger:         logger,
		snippets:       &models.SnippetModel{Pool: db},
		users:          &models.UserModel{Pool: db},
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}
	// tlsConfig := &tls.Config{
	// 	CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	// }

	srv := http.Server{
		Addr:     fmt.Sprint(":", utils.GetEnv("PORT", "4000")),
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
		// TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info(fmt.Sprintf("Creating Server on http://localhost:%v", utils.GetEnv("PORT", "4000")))
	fmt.Println()
	// err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err = srv.ListenAndServe()
	app.logger.Error(err.Error())

	defer db.Close()
}
