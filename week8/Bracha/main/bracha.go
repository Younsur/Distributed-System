package main

import (
	"fmt"
)

const (
	f = 1
	N = 3*f + 1
)

type Bracha struct {
	Parties []*Party
}

func NewBracha() *Bracha {
	bracha := &Bracha{}
	for i := 0; i < N; i++ {
		node := &Party{
			IsHonesty: true,
			Id:        i,
			Echo:      make(map[string]int),
			Ready:     make(map[string]int),
			Bracha:    bracha,
			OutputCh:  make(chan string),
			IsLocked:  false,
			Content:   "",
		}
		bracha.Parties = append(bracha.Parties, node)
	}
	return bracha
}

func main() {
	bracha := NewBracha()
	bracha.Parties[0].IsHonesty = true
	bracha.Parties[1].IsHonesty = true
	bracha.Parties[2].IsHonesty = true

	bracha.Parties[0].Broadcast(Message{Type: VALUE, Content: "GenshinStart"}, bracha)
	for _, party := range bracha.Parties {
		fmt.Printf("party %d output %s\n", party.Id, party.Content)
	}
}
