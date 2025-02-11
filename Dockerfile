# Build stage
FROM golang:1.23.3-alpine3.20 AS build
WORKDIR /app

COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.20 
WORKDIR /app
COPY --from=build /app/main .
COPY .env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]