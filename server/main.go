package main

import (
	"log"
	"net/http"
)

func main() {
	server := newServer()

	go server.listen()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleIO(server, w, r)
	})
	log.Fatal(http.ListenAndServe(":8123", nil))
}
