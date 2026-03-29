go-migrate 需要手动install 才能支持对应数据库
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest


sqlc
sqlc init // 初始化
document:
https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html


migrate create -ext sql  -dir db/migration -seq add_users

docker exec -it postgres psql -U root -d simple_bank

begin;
select * from accounts where id=1 for update;
update accounts set balance=1000 where id=1;
commit;

外键约束+lock

go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc