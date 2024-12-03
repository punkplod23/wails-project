run:
	wails dev

build:
	go env -w CGO_ENABLED=0
	wails build



lintgo:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

