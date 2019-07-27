.PHONY: test test-unit test-e2e test-middlewares test-parameters imports fmt lint install build build-darwin build-linux build-windows

test: test-unit test-e2e

test-unit: imports pkg/generator/templates.go
	@(go list ./... | grep -v "vendor/" | grep -v "e2e" | xargs -n1 go test -cover)

test-e2e: test-middlewares test-parameters

test-middlewares: imports pkg/generator/templates.go internal/test/e2e/go-oas-server-test
	@(internal/test/e2e/go-oas-server-test generate --file internal/test/e2e/middlewares/api.json --output-dir=internal/test/e2e/middlewares)
	@(go test -v -cover github.com/michaelsauter/go-oas-server/internal/test/e2e/middlewares)

test-parameters: imports pkg/generator/templates.go internal/test/e2e/go-oas-server-test
	@(internal/test/e2e/go-oas-server-test generate --file internal/test/e2e/parameters/api.json --output-dir=internal/test/e2e/parameters)
	@(go test -v -cover github.com/michaelsauter/go-oas-server/internal/test/e2e/parameters)

imports:
	@(goimports -w .)

fmt:
	@(gofmt -w .)

lint:
	@(go mod download && golangci-lint run)

install: imports
	@(go install)

build: imports build-linux build-darwin build-windows pkg/generator/templates.go

build-linux: imports pkg/generator/templates.go
	cd cmd/go-oas-server && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o go-oas-server-linux-amd64

build-darwin: imports pkg/generator/templates.go
	cd cmd/go-oas-server && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o go-oas-server-darwin-amd64

build-windows: imports pkg/generator/templates.go
	cd cmd/go-oas-server && GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o go-oas-server-windows-amd64.exe

pkg/generator/templates.go: pkg/templates/*
	@(echo "Generating pkg/generator/templates.go")
	@(./generate_templates_go.sh)

internal/test/e2e/go-oas-server-test: cmd/go-oas-server/main.go go.mod go.sum pkg/commands/* pkg/generator/*
	@(echo "Generating E2E test binary")
	@(cd cmd/go-oas-server && go build -o ../../internal/test/e2e/go-oas-server-test)
