package main

import (
	"log"
	"os"

	"github.com/akross26/final-project/pkg/db"
	"github.com/akross26/final-project/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		return err
	}
	defer db.Close()

	return server.StartServer(port)
}
