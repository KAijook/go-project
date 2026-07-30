
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test migrateup1 migratedown1 server


postgres:
	docker run --name my-postgres -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:17


createdb:
	docker exec -it my-postgres createdb --username=postgres --owner=postgres simple_bank


dropdb:
	docker exec -it my-postgres dropdb --username=postgres simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go