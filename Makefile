createdb:
	docker exec -it postggres12 createdb --username=root --owner=root simple_bank
postgres:
	docker exec -it postggres12 psql -U postgres -d simple_bank
dropdb:
	docker exec -it postggres12 dropdb simple_bank

createmigration:
	 migrate create -ext sql -dir ./db/migratuon -seq init_schema


migrateup:
	migrate -path ./db/migration -database "postgresql://postgres:iloveyou044@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path ./db/migration -database "postgresql://postgres:iloveyou044@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path ./db/migration -database "postgresql://postgres:iloveyou044@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate
test:
	go test -v ./db/sqlc
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go  github.com/techschool/simplebank/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock