createMigration: 
	migrate create -ext sql -dir internal/db/migrations -seq name 

postgres:
	docker run --name bank -e POSTGRES_USER=dev -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=simple_bank -p 5433:5432 -d postgres:17-alpine

createdb:
	docker exec -it bank createdb username=dev --owner=dev simple_bank

dropdb:
	docker exec -it bank dropdb simple_bank

migrateup:
	migrate --path internal/db/migrations -database "postgresql://dev:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up	

migratedown:
	migrate --path internal/db/migrations -database "postgresql://dev:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down	

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go	

.PHONY:createMigration postgres createdb dropdb migrateup migratedown sqlc test server
