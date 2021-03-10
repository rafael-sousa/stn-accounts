.PHONY: start
start:
	docker-compose -f docker-compose.yml up -d

.PHONY: stop
stop:
	docker-compose -f docker-compose.yml down

.PHONY: lint
lint:
	docker-compose -f docker-compose.yml down

.PHONY: lint
lint:
	go mod tidy
    go vet ./...
    go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: test
test:
	go build -o bin/main.exe cmd/stn-accounts/main.go
	docker-compose build

.PHONY: run
run:
	go run cmd/stn-accounts/main.go


.PHONY: gen-src
gen-src:
	swag init -g cmd/stn-accounts/main.go
