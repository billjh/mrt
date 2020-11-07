package main

import (
	"log"
	"net/http"
)

func main() {
	navigator := NewNavigator()

	http.HandleFunc("/api/navigate/v1", navigator.handleV1)
	http.HandleFunc("/api/navigate/v2", navigator.handleV2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
