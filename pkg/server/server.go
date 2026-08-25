package server

import (
	"log"
	"net/http"

	"github.com/akross26/final-project/pkg/api"
)

func StartServer(port string) error {
	const webDir = "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init()

	addr := ":" + port

	log.Printf("Сервер запущен на http://localhost%s", addr)

	return http.ListenAndServe(addr, nil)
}
