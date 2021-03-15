package main

import (
	"log"
	"net/http"
	"time"

	"/Users/min/go/src/learn_go_with_test/http-server"

	"github.com/gomodule/redigo/redis"
)

func main() {
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	server := NewPlayerServer(NewRedisPlayerStore("go_server_score"))

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
