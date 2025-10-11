migration: 
	migrate create -ext sql -dir internal/db/migrations -seq name 

postgres:
	docker run --name bank -e POSTGRES_USER=dev -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=simple_bank -p 5433:5432 -d postgres:17-alpine

createdb:
	docker exec -it bank createdb username=dev --owner=dev simple_bank

dropdb:
	docker exec -it bank dropdb simple_bank

migrateup:
	migrate --path internal/db/migrations -database "postgresql://dev:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up	

migrateup1:
	migrate --path internal/db/migrations -database "postgresql://dev:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up 1	

migratedown:
	migrate --path internal/db/migrations -database "postgresql://dev:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down	

migratedown1:
	migrate --path internal/db/migrations -database "postgresql://dev:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination internal/db/mock/store.go -source internal/db/sqlc/store.go -aux_files github.com/Doris-Mwito5/simple-bank/internal/db/sqlc=internal/db/sqlc/querier.go

forcereset:
	migrate -path internal/db/migrations -database "postgresql://dev:secret@localhost:5433/simple_bank?sslmode=disable" force 2


.PHONY:createMigration postgres createdb dropdb migrateup1 migratedown migratedown1 sqlc test server mock forcereset
