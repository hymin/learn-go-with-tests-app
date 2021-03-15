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

func (i *InMemoryPlayerStore) GetLeague() (league []Player) {
	for k, v := range i.scores {
		league = append(league, Player{k, v})
	}
	return
}
