package main

import (
	"fmt"
	"log"
	"net/http"
)

// redirectHandler takes an HTTP request and redirects to a special place
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.youtube.com/embed/dQw4w9WgXcQ", 302)
}

func main() {
	fmt.Printf("Hello, world.\n")

	http.HandleFunc("/", redirectHandler)
	log.Fatal(http.ListenAndServe(":1441", nil))
}
