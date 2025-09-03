package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq" // PostgreSQL driver
    "github.com/techschool/simplebank/api"
    db "github.com/techschool/simplebank/db/sqlc"
)

const (
    dbDriver      = "postgres"
    dbSource      = "postgresql://postgres:iloveyou044@localhost:5432/simple_bank?sslmode=disable"
    serverAddress = "0.0.0.0:"
    

)

func main() {
    // 1. Kết nối database
    conn, err := sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("cannot connect to db: ", err)
    }

    // 2. Tạo store từ sqlc (Queries)
    store := db.NewStore(conn)

    // 3. Tạo server
    server := api.NewServer(store)

    // 4. Start server
    err = server.Start(serverAddress)
    if err != nil {
        log.Fatal("cannot start server: ", err)
    }
}
