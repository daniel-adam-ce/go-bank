include .env

postgres:
	docker run --name postgres17 --network bank-network -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:17-alpine

postgresstart:
	docker start postgres17

createdb:
	docker exec -it postgres17 createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_NAME}

dropdb:
	docker exec -t postgres17 dropdb ${POSTGRES_NAME}

migrateup:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_NAME}${PG_OPTIONS}" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_NAME}${PG_OPTIONS}" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_NAME}${PG_OPTIONS}" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_NAME}${PG_OPTIONS}" -verbose down 1

migratecreate: 
	migrate create -ext sql -dir db/migration -seq ${name}

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

testpkg:
	go test ./...  -coverpkg=./... -coverprofile ./coverage.out
	go tool cover -func ./coverage.out

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/daniel-adam-ce/go-bank/db/sqlc Store

dbdocs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=go_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

redis:
	docker run --name redis -p 6379:6379 -d redis:7.4.1-alpine

redisstart:
	docker start redis

.PHONY: postgres postgresstart createdb dropdb migrateup migratedow nmigrateup1 migratedown1 sqlc server mock dbdocs db_schema test testpkg proto redis redisstart