package main

import (
	"log"
	"net/http"

	"github.com/pakohan/go-faster/api/controller"

	_ "github.com/pakohan/go-faster/api/controller/logs"
)

func main() {
	controller.Init()

	server := &http.Server{
		Addr:    ":8080",
		Handler: controller.Container,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
