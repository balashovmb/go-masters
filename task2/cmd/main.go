package main

import (
	"flag"
	"log"

	"go-masters/task2/internal/api"
	"net/http"
)

func main() {
	var port string
	flag.StringVar(&port, "p", "8080", "keyword")
	flag.Parse()
	http.HandleFunc("/rate", api.Rate)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
