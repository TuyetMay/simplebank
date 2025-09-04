package main

import (
    "database/sql"
    "log"
    "github.com/techschool/simplebank/util"

    _ "github.com/lib/pq" // PostgreSQL driver
    "github.com/techschool/simplebank/api"
    db "github.com/techschool/simplebank/db/sqlc"
)

func main() {

    config, err := util.LoadConfig(".")
    if err != nil {
        log.Fatal("cannot load config:", err)
    }


    // 1. Kết nối database
    conn, err := sql.Open(config.DBDriver, config.DBSource)
    if err != nil {
        log.Fatal("cannot connect to db: ", err)
    }

    // Test connection
    if err = conn.Ping(); err != nil {
        log.Fatal("cannot ping db: ", err)
    }

    // 2. Tạo store từ sqlc (Queries)
    store := db.NewStore(conn)

    // 3. Tạo server
    server := api.NewServer(store)

    // 4. Start server
    log.Printf("Starting server on %s", config.ServerAddress)
    err = server.Start(config.ServerAddress)
    if err != nil {
        log.Fatal("cannot start server: ", err)
    }
}
