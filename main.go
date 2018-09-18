package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
