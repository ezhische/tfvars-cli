test:
	go test ./... -v
build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=dev-build" -o ./tfvar-cli ./cmd/main.go
lint-install:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin
lint:
	golangci-lint run ./... -v
