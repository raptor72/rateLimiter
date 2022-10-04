BINARY_NAME=rateLimiter

.PHONY: build
build:
	go build -o ${BINARY_NAME} -mod=vendor main.go

run:
	docker-compose -f service-limiter.yml up

run-background:
	docker-compose -f service-limiter.yml up -d

test: run-background
	go test -v ./... & docker-compose -f service-limiter.yml down

test-race: run-background
	go test -race -count 100 ./... & docker-compose -f service-limiter.yml down

lint:
	golangci-lint run
