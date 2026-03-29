DB_URL=postgresql://root:secret@127.0.0.1:5432/simple_bank?sslmode=disable

network:
	docker network create bank-network

network-connect:
	docker network disconnect bank-network simplebank
	docker network disconnect bank-network postgres
	docker network connect bank-network simplebank
	docker network connect bank-network postgres

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

dev:
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

build:
	docker rmi simplebank:v1.0
	docker build -t simplebank:v1.0 .

run:
	docker run --name simplebank --network bank-network -p 8081:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:secret@postgres:5432/simple_bank?sslmode=disable" simplebank:v1.0

proto:
	del /f /q pb\*.go
	del /f /q doc\swagger\*swagger.json
	protoc --proto_path=./proto --go_out=./pb --go_opt=paths=source_relative \
		--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out ./pb --grpc-gateway_opt paths=source_relative \
		--openapiv2_out=./doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
		proto/*.proto
	statik -src=./doc/swagger -dest=./doc/

evans:
	evans --host localhost --port 9091 -r repl
	package pb
	service SimpleBank

.PHONY: migrateup1 migratedown1 sqlc test mock createdb postgres network dropdb clean mock dev run proto evans