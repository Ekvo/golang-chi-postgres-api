#!/bin/bash

# database
go get github.com/lib/pq
go get github.com/joho/godotenv

# transport
go get github.com/go-chi/chi/v5

go mod tidy

# run app
go run cmd/app/main.go
