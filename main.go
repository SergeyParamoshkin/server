package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	var (
		router *mux.Router
		port   string
		err    error
	)
	t := telnet
	t.Check()
	port = ":8080"
	router = NewRouter()

	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
