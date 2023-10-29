package main

import (
	"fmt"
	"log"
	"sync"
)

const (
	VALUE string = "value"
	ECHO  string = "echo"
	READY string = "ready"
)

type Party struct {
	Id        int
	SendReady bool
	Mu        sync.Mutex
	Echo      map[string]int
	Ready     map[string]int
	OutputCh  chan string
	IsHonesty bool
	Bracha    *Bracha
	IsLocked  bool
	Content   string
}

type Message struct {
	Type    string
	Content string
}

func (p *Party) ReceiveMessage(m Message) {
	p.Mu.Lock()
	p.IsLocked = true
	var broadcastMessage *Message = nil
	switch m.Type {
	case VALUE:
		if p.IsHonesty {
			broadcastMessage = &Message{Type: ECHO, Content: m.Content}
		} else {
			broadcastMessage = &Message{Type: ECHO, Content: m.Content + fmt.Sprint("fake")}
		}
	case ECHO:
		p.Echo[m.Content]++
		log.Printf("party %d receives echo %s total %d", p.Id, m.Content, p.Echo[m.Content])
		if p.Echo[m.Content] >= 2*f+1 && !p.SendReady {
			p.SendReady = true
			broadcastMessage = &Message{Type: READY, Content: m.Content}
		}
	case READY:
		p.Ready[m.Content]++
		if p.Ready[m.Content] >= 2*f+1 {
			p.Content = m.Content
		} else if p.Ready[m.Content] >= f+1 && !p.SendReady {
			p.SendReady = true
			broadcastMessage = &Message{Type: READY, Content: m.Content}
		}
	}
	p.IsLocked = false
	p.Mu.Unlock()
	if broadcastMessage != nil {
		p.Broadcast(*broadcastMessage, p.Bracha)
	}
}

func (p *Party) Broadcast(m Message, bracha *Bracha) {
	for i := 0; i < N; i++ {
		if p.IsLocked {
			p.IsLocked = false
			p.Mu.Unlock()
		}
		if bracha.Parties[i].IsLocked {
			bracha.Parties[i].IsLocked = false
			bracha.Parties[i].Mu.Unlock()
		}
		bracha.Parties[i].ReceiveMessage(m)
	}
}
