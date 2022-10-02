BINARY_NAME=rateLimiter

.PHONY: build
build:
	go build -o ${BINARY_NAME} -mod=vendor main.go

run:
	docker-compose -f service-limiter.yml up -d

test: run
	go test -v ./... & docker-compose -f service-limiter.yml down

