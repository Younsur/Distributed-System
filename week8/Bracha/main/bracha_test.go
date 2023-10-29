package main

import (
	"fmt"
	"testing"
)

const (
	f = 1
	N = 3*f + 1
)

func TestBracha(t *testing.T) {
	bracha := NewBracha()
	bracha.Parties[0].IsHonesty = true
	bracha.Parties[1].IsHonesty = true
	bracha.Parties[2].IsHonesty = true

	bracha.Parties[0].Broadcast(Message{Type: VALUE, Content: "GenshinStart"}, bracha)
	for _, party := range bracha.Parties {
		fmt.Printf("party %d output %s\n", party.Id, party.Content)
	}
}
