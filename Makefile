.PHONY: run

run:
	go run cmd/main.go

keygen:
	go run tools/rsa/keygen.go

# make apigen ARGS="sample"
apigen:
	go run tools/apigen.go $(ARGS)