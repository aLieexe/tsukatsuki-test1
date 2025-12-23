hello:
	echo "Hello"

dev:
	air

fmt:
	go fmt ./src/...

test:
	go test -v ./cmd/web -count=1

test-internal:
	go test -v ./internal/models -count=1

coverage:
	mkdir -p tmp
	go test -covermode=count -coverprofile=tmp/coverage.out ./...
	go tool cover -func=tmp/coverage.out
	go tool cover -html=tmp/coverage.out -o tmp/coverage.html

debug: 
	go run ./cmd/web -debug