package main

import (
    "database/sql"
    "log"
    "github.com/techschool/simplebank/util"

    _ "github.com/lib/pq" // PostgreSQL driver
    "github.com/techschool/simplebank/api"
    db "github.com/techschool/simplebank/db/sqlc"
)

const (
    dbDriver      = "postgres"
    dbSource      = "postgresql://user_name:pass_word@localhost:5432/simple_bank?sslmode=disable"
    serverAddress = "0.0.0.0:"
    

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

    // 2. Tạo store từ sqlc (Queries)
    store := db.NewStore(conn)

    // 3. Tạo server
    server := api.NewServer(store)

    // 4. Start server
    err = server.Start(config.ServerAddress)
    if err != nil {
        log.Fatal("cannot start server: ", err)
    }
}
