createdb:
	docker exec -it postggres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postggres12 dropdb simple_bank

createmigration:
	 migrate create -ext sql -dir ./db/migratuon -seq init_schema

migrateup:
	migrate -path ./db/migration -database "postgresql://postgres:iloveyou044@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migration -database "postgresql://postgres:iloveyou044@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate
test:
	go test -v ./db/sqlc

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test