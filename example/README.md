## Migrate

```
go install github.com/rubenv/sql-migrate/...@latest
sql-migrate up
```

## Generate

```
# sqlc
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
sqlc generate

# gqlgen
go run github.com/99designs/gqlgen generate

# gRPC
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/*.proto

# mapc
go run .
```
