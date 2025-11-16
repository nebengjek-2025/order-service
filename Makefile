.PHONY: install test cover run build start clean

install:
	@go mod download

test:
	@echo "Run unit testing ..."
	@mkdir -p ./coverage && \
	go test -v ./test/ -coverprofile=./coverage/coverage.out -covermode=atomic ./...

cover: test
	@echo "Generating coverprofile ..."
	@go tool cover -func=./coverage/coverage.out && \
	go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html

run:
	@go run ./src/cmd/app/main.go

run-worker:
	@go run ./cmd/worker/main.go

build:
	@go build -tags musl -o main ./cmd/web

start:
	@./main

clean:
	@echo "Cleansing the last built ..."
	@rm -rf bin

migrate:
	@echo "Running database migration ..."
	@migrate -database "mysql://root:my-secret-pw@tcp(localhost:3306)/nebengjek?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up