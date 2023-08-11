.PHONY: mod-download-local
mod-download-local:
	go mod download && go mod tidy

.PHONY: build-local
build-local:
	GODEBUG="tarinsecurepath=0,zipinsecurepath=0" go build -v

.PHONY: test-unit-with-coverage
test-with-coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic -v ./tests/units/... -coverpkg=./...

.PHONY: test-e2e
test-e2e:
	rm -rf e2e-test-results
	go build -v -o /tmp/kermoo .
	KERMOO_BINARY="/tmp/kermoo" go test -v ./tests/e2e/...
