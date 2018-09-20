package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/filter", Filter)
	router.POST("/prioritize", Prioritize)

	log.Fatal(http.ListenAndServe(":8888", router))
}
