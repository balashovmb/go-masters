package main

import (
	"flag"
	"go-masters/task2/internal/api"
	"log"
	"net/http"
)

func main() {
	api := &api.API{}
	var port string
	flag.StringVar(&port, "p", "8080", "keyword")
	flag.Parse()
	http.HandleFunc("/rate", api.Rate)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
