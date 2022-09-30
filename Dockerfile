FROM golang:1.18-alpine AS build
ENV LANG C.UTF-8
WORKDIR /rateLimiter
COPY . /rateLimiter

RUN go build -o rateLimiter -mod=vendor main.go


FROM alpine:3.9
COPY --from=build /rateLimiter .
EXPOSE 8080

CMD ./rateLimiter