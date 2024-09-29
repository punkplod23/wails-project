run:
	wails dev

build:
	wails build

lintgo:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

