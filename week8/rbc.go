package aab

import (
	"Dory/internal/party"
	"Dory/pkg/protobuf"
	"encoding/binary"
	"sync"
)

func RBCMain(p *party.HonestParty, outputChannel chan []byte) {
	index := uint32(0)

	syncChan := make(map[uint32]chan bool)
	syncChan[0] = make(chan bool, 1)
	syncChan[0] <- true

	for {
		index++
		id := make([]byte, 4)
		binary.BigEndian.PutUint32(id, index)
		syncChan[index] = make(chan bool, 1)
		data := <-p.GetMessage("PROPOSE", []byte{0b11111111})

		go RBC(p, index, id, data.Data, outputChannel, syncChan)
	}
}

func RBC(p *party.HonestParty, index uint32, id []byte, data []byte, outputChannel chan []byte, syncChan map[uint32]chan bool) {
	sendReady := false //是否已发送Ready
	echo := make(map[string]int)
	ready := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	if p.PID >= 2*p.F+1 {
		data[0] = data[0] + 6
		err := p.Broadcast(&protobuf.Message{
			Data: data,
			Type: "ECHO",
			Id:   id,
		})
		if err != nil {
			return
		}
	} else {
		err := p.Broadcast(&protobuf.Message{
			Data: data,
			Type: "ECHO",
			Id:   id,
		})
		if err != nil {
			return
		}
	}

	go Echo(p, id, echo, &sendReady, &mu, &wg)
	go Ready(p, index, id, ready, &sendReady, outputChannel, &mu, &wg, syncChan)

	wg.Wait()
}

func Echo(p *party.HonestParty, id []byte, echo map[string]int, sendReady *bool, mu *sync.Mutex, wg *sync.WaitGroup) {
	for {
		f := int(p.F)
		m := <-p.GetMessage("ECHO", id)
		mu.Lock()
		if !*sendReady {
			if m != nil {
				k := string(m.Data)
				v, ok := echo[k]
				if ok {
					echo[k] = v + 1
				} else {
					echo[k] = 1
				}
				if echo[k] >= 2*f+1 {
					err := p.Broadcast(&protobuf.Message{
						Data: m.Data,
						Type: "READY",
						Id:   id,
					})
					if err != nil {
						return
					}
					*sendReady = true
					break
				}
			}
		} else {
			break
		}
		mu.Unlock()
	}
	mu.Unlock()
	wg.Done()
}

func Ready(p *party.HonestParty, index uint32, id []byte, ready map[string]int, sendReady *bool, outputChannel chan []byte, mu *sync.Mutex, wg *sync.WaitGroup, syncChan map[uint32]chan bool) {
	for {
		f := int(p.F)
		m := <-p.GetMessage("READY", id)
		mu.Lock()
		if m != nil {
			k := string(m.Data)
			v, ok := ready[k]
			if ok {
				ready[k] = v + 1
			} else {
				ready[k] = 1
			}
			if ready[k] >= f+1 && !*sendReady {
				err := p.Broadcast(&protobuf.Message{
					Data: m.Data,
					Type: "READY",
					Id:   id,
				})
				if err != nil {
					return
				}
				*sendReady = true
				mu.Unlock()
				continue
			}
			if ready[k] >= 2*f+1 {
				<-syncChan[index-1]
				outputChannel <- m.Data
				syncChan[index] <- true
				mu.Unlock()
				break
			}
		}
		mu.Unlock()
	}
	wg.Done()
}
