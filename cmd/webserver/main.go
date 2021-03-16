package main

import (
	"log"
	"net/http"
	"time"

	poker "github.com/hymin/learn-go-with-tests-app"

	"github.com/gomodule/redigo/redis"
)

func main() {
	//test
	poker.pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	server := poker.NewPlayerServer(poker.NewRedisPlayerStore("go_server_score"))

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
