DB_URL=sqlite3://sqlite.db

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up
migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover -short ./...
server:
	go run main.go

.PHONY: migrateup  migratedown sqlc test server