.PHONY: run

run:
	go run cmd/main.go

keygen:
	go run tools/rsa/keygen.go