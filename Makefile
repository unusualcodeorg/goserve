.PHONY: run

install:
	go mod tidy

build:
	go build -o build/server cmd/main.go

run:
	go run cmd/main.go

test:
	go test -v ./...

cover:
	go test -cover ./...

keygen:
	go run .tools/rsa/keygen.go

# make apigen ARGS="sample"
apigen:
	go run .tools/apigen.go $(ARGS)