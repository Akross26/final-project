package main

import (
	"log"
	"os"

	"github.com/akross26/final-project/pkg/db"
	"github.com/akross26/final-project/pkg/server"

	_ "modernc.org/sqlite"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := server.StartServer(port); err != nil {
		log.Fatal(err)
	}
}
