postgres:
	docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root api_db

dropdb:
	docker exec -it postgres15 dropdb api_db

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/api_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/api_db?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/JekaTka/services-api/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock
