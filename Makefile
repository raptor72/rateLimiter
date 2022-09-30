
build-app:
	go build -o rateLimiter -mod=vendor main.go

db-up:
	docker-compose -f build/services-local.yml up --build -d

db-down:
	docker-compose -f build/services-local.yml down

app-run: build-app
	./rateLimiter
