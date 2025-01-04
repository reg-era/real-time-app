package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := http.NewServeMux()

	log.Printf("Route server running on http://localhost:%s\n", 8080)
	log.Fatalln(http.ListenAndServe(":8080", router))
}
