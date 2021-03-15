package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
)

func TestRecordingInMemory(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, want, got)
	})
}

func TestRecordingRedis(t *testing.T) {
	redisKey := "test_score_hash"
	pool = newRedisPool(redisKey)
	store := NewRedisPlayerStore(redisKey)
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(t, player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, want, got)
	})
}

func newRedisPool(redisKey string) *redis.Pool {
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	conn := pool.Get()
	conn.Do("del", redisKey)
	return pool
}
