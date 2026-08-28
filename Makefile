.PHONY: swag

swag:
	swag init -g server.go -o docs

run:
	go run server.go