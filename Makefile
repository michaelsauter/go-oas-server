.PHONY: test test-unit test-e2e test-middlewares test-parameters imports fmt lint install build build-darwin build-linux build-windows

test: test-unit test-e2e

test-unit: imports generator/templates.go
	@(go list ./... | grep -v "vendor/" | grep -v "e2e" | xargs -n1 go test -cover)

test-e2e: test-middlewares test-parameters

test-middlewares: imports generator/templates.go test/e2e/go-oas-server-test
	@(test/e2e/go-oas-server-test generate --file test/e2e/middlewares/api.json --output-dir=test/e2e/middlewares)
	@(go test -v -cover github.com/michaelsauter/go-oas-server/test/e2e/middlewares)

test-parameters: imports generator/templates.go test/e2e/go-oas-server-test
	@(test/e2e/go-oas-server-test generate --file test/e2e/parameters/api.json --output-dir=test/e2e/parameters)
	@(go test -v -cover github.com/michaelsauter/go-oas-server/test/e2e/parameters)

imports:
	@(goimports -w .)

fmt:
	@(gofmt -w .)

lint:
	@(go mod download && golangci-lint run)

install: imports
	@(go install)

build: imports build-linux build-darwin build-windows generator/templates.go

build-linux: imports generator/templates.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o go-oas-server-linux-amd64 -v github.com/michaelsauter/go-oas-server

build-darwin: imports generator/templates.go
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o go-oas-server-darwin-amd64 -v github.com/michaelsauter/go-oas-server

build-windows: imports generator/templates.go
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o go-oas-server-windows-amd64.exe -v github.com/michaelsauter/go-oas-server

generator/templates.go: templates/*
	@(echo "Generating generator/templates.go")
	@(./generate_templates_go.sh)

test/e2e/go-oas-server-test: main.go go.mod go.sum commands/* generator/*
	@(echo "Generating E2E test binary")
	@(go build -o test/e2e/go-oas-server-test)
