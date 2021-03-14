package main

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
