package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)
	http.HandleFunc("/share", shareHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
