package main

import (
	"log"
	"net/http"

	poker "github.com/hymin/learn-go-with-tests-app"
)

func main() {
	server := poker.NewPlayerServer(poker.NewRedisPlayerStore("go_server_score"))

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
