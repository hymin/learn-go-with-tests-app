package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	scores map[string]int
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) (int, bool) {
	score, ok := i.scores[name]
	return score, ok
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.scores[name]++
}

func NewRedisPlayerStore(redisKey string) *RedisPlayerStore {
	return &RedisPlayerStore{redisKey}
}

type RedisPlayerStore struct {
	scoreName string
}

func (r *RedisPlayerStore) GetPlayerScore(name string) (int, bool) {
	conn := pool.Get()
	defer conn.Close()
	score, err := redis.Int(conn.Do("hget", r.scoreName, name))
	ok := err == nil
	return score, ok
}

func (r *RedisPlayerStore) RecordWin(name string) {
	conn := pool.Get()
	defer conn.Close()
	_, err := redis.Int(conn.Do("hincrby", r.scoreName, name, 1))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	server := &PlayerServer{NewRedisPlayerStore("go_server_score")}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
