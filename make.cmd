@echo off
IF "%1"=="start" docker-compose -f docker-compose.yml up -d
IF "%1"=="stop" docker-compose -f docker-compose.yml down
IF "%1"=="lint" (
    go mod tidy
    go vet ./...
    go fmt ./...
    golint ./...
)
IF "%1"=="test" go test -v ./...
IF "%1"=="build" (
    go build -o bin/main.exe cmd/stn-accounts/main.go
    docker-compose build
)
IF "%1"=="run" go run cmd/stn-accounts/main.go
IF "%1"=="gen-src" swag init -g cmd/stn-accounts/main.go

