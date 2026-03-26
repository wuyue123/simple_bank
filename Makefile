DB_URL=postgresql://root:secret@127.0.0.1:5432/simple_bank?sslmode=disable
network:
	docker network create bank-network

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

run:
	go run main.go

createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb --username=root simple_bank

mock:
	mockgen -package mockdb -destination db/mock/store.go simpleblank/db/sqlc Store
#	mockgen -package mockwk -destination worker/mock/distributor.go github.com/techschool/simplebank/worker TaskDistributor


migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up


migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down


migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1


sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

clean:
	go clean -modcache
	go clean -cache

mock:
	mockgen -destination db/mock/store.go -package mockdb pxsemic.com/simplebank/db/sqlc Store

.PHONY: migrateup1 migratedown1 sqlc test mock createdb postgres network dropdb clean mock