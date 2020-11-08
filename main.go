package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	navigator := NewNavigator()
	httpPort := ":8080"

	http.HandleFunc("/api/navigate/v1", navigator.handleV1)
	http.HandleFunc("/api/navigate/v2", navigator.handleV2)

	fmt.Printf("Listening on %s\n", httpPort)
	log.Fatal(http.ListenAndServe(httpPort, nil))
}
