package poker

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool = &redis.Pool{
	MaxIdle:     10,
	IdleTimeout: 240 * time.Second,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
}

func NewRedisPlayerStore(keyScore string) *RedisPlayerStore {
	return &RedisPlayerStore{keyScore}
}

type RedisPlayerStore struct {
	keyScore string
}

func (r *RedisPlayerStore) GetPlayerScore(name string) (int, bool) {
	conn := pool.Get()
	defer conn.Close()
	score, err := redis.Int(conn.Do("hget", r.keyScore, name))
	ok := err == nil
	return score, ok
}

func (r *RedisPlayerStore) RecordWin(name string) {
	conn := pool.Get()
	defer conn.Close()
	_, err := redis.Int(conn.Do("hincrby", r.keyScore, name, 1))
	if err != nil {
		fmt.Println(err)
	}
}

func (r *RedisPlayerStore) GetLeague() (league []Player) {
	conn := pool.Get()
	defer conn.Close()
	content, err := redis.StringMap(conn.Do("hgetall", r.keyScore))
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range content {
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Println(err)
		}
		league = append(league, Player{k, i})
	}
	return
}
