include .env

postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_NAME}

dropdb:
	docker exec -t postgres17 dropdb ${POSTGRES_NAME}

migrateup:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_NAME}?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_NAME}?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/daniel-adam-ce/go-bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server mock