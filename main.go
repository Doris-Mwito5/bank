package main

import (
	"database/sql"
	"github.com/Doris-Mwito5/simple-bank/api"
	db "github.com/Doris-Mwito5/simple-bank/internal/db/sqlc"
	"github.com/Doris-Mwito5/simple-bank/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal("cannot start server: %w", err)
	}	
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
