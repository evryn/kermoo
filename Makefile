.PHONY: mod-download-local
mod-download-local:
	go mod download && go mod tidy

.PHONY: build-local
build-local:
	GODEBUG="tarinsecurepath=0,zipinsecurepath=0" go build -v `go list ./... | grep -v 'test/e2e'`

.PHONY: test-with-coverage
test-with-coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic
