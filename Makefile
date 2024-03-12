default: build

test:
	go mod download && go test ./...

e2e: install
	cd integration && tflint --chdir=.

build:
	go build

install:
	go run install/main.go

lint:
	golint --set_exit_status $$(go list ./...)
	go vet ./...

tools:
	go install golang.org/x/lint/golint@latest

.PHONY: test e2e build install lint tools
