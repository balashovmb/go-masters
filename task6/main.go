package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {
	var port string
	flag.StringVar(&port, "p", "8080", "port")
	flag.Parse()
	http.HandleFunc("/fibo", fibo)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func fibo(w http.ResponseWriter, r *http.Request) {
	n := r.URL.Query().Get("N")
	if n == "" {
		http.Error(w, "N is empty", http.StatusBadRequest)
		return
	}
	nInt, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(strconv.Itoa(calcFibo(nInt))))
}

func calcFibo(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	prevFib := 0
	currFib := 1

	for i := 2; i < n; i++ {
		nextFib := prevFib + currFib
		prevFib = currFib
		currFib = nextFib
	}

	return currFib
}
