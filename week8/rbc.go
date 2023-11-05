package aab

import (
	"Dory/internal/party"
	"Dory/pkg/protobuf"
	"encoding/binary"
	"fmt"
	"sync"
)

func RBCMain(p *party.HonestParty, inputChannel chan []byte, outputChannel chan []byte) {
	index := uint32(0)
	mm := make(map[uint32]chan bool)
	mm[0] = make(chan bool, 1)
	mm[0] <- true

	for {
		index++
		id := make([]byte, 4)
		binary.BigEndian.PutUint32(id, index)
		mm[index] = make(chan bool, 1)
		data := <-p.GetMessage("PROPOSE", []byte{0b11111111})

		if data.Data[0] != id[0] {
			fmt.Errorf("ERROR")
		}

		go RBC(p, index, id, data.Data, outputChannel, mm)
	}
}

func RBC(p *party.HonestParty, index uint32, id []byte, data []byte, outputChannel chan []byte, mm map[uint32]chan bool) {
	sendReady := false
	ready1 := make(map[string]int)
	ready2 := make(map[string]int)
	mu := new(sync.Mutex)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for {
			f := int(p.F)
			m := <-p.GetMessage("ECHO", id)
			mu.Lock()
			if !sendReady {
				if m != nil {
					k := string(m.Data)
					v, ok := ready1[k]
					if ok {
						ready1[k] = v + 1
					} else {
						ready1[k] = 1
					}
					if ready1[k] >= 2*f+1 {
						p.Broadcast(&protobuf.Message{
							Data: m.Data,
							Type: "READY",
							Id:   id,
						})
						sendReady = true
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
	}()

	go func() {
		for {
			f := int(p.F)
			m := <-p.GetMessage("READY", id)
			mu.Lock()
			if m != nil {
				k := string(m.Data)
				v, ok := ready2[k]
				if ok {
					ready2[k] = v + 1
				} else {
					ready2[k] = 1
				}
				if ready2[k] >= f+1 && !sendReady {
					p.Broadcast(&protobuf.Message{
						Data: m.Data,
						Type: "READY",
						Id:   id,
					})
					sendReady = true
					mu.Unlock()
					continue
				}
				if ready2[k] >= 2*f+1 {
					<-mm[index-1]
					outputChannel <- m.Data
					mm[index] <- true
					mu.Unlock()
					break
				}
			}
			mu.Unlock()
		}
		wg.Done()
	}()
	//go Echo(p, id, data, ready1, sendReady, mu, &wg)
	//go Ready(p, index, id, data, ready2, sendReady, outputChannel, mm, mu, &wg)

	if p.PID > 2*p.F {
		data[0] = data[0] + 1
		p.Broadcast(&protobuf.Message{
			Data: data,
			Type: "ECHO",
			Id:   id,
		})
	} else {
		p.Broadcast(&protobuf.Message{
			Data: data,
			Type: "ECHO",
			Id:   id,
		})
	}

	wg.Wait()
}

/*func Echo(p *party.HonestParty, id []byte, data []byte, ready map[string]int, sendReady bool, mu *sync.Mutex, wg *sync.WaitGroup) {
	for {
		f := int(p.F)
		m := <-p.GetMessage("ECHO", id)
		mu.Lock()
		if !sendReady {
			if m != nil {
				k := string(m.Data)
				v, ok := ready[k]
				if ok {
					ready[k] = v + 1
				} else {
					ready[k] = 1
				}
				if ready[k] >= 2*f+1 {
					p.Broadcast(&protobuf.Message{
						Data: m.Data,
						Type: "READY",
						Id:   id,
					})
					sendReady = true
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
}*/

/*func Ready(p *party.HonestParty, index uint32, id []byte, data []byte, ready map[string]int, sendReady bool, outputChannel chan []byte, mm map[uint32]chan bool, mu *sync.Mutex, wg *sync.WaitGroup) {
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
			if ready[k] >= f+1 && !sendReady {
				p.Broadcast(&protobuf.Message{
					Data: m.Data,
					Type: "READY",
					Id:   id,
				})
				sendReady = true
				mu.Unlock()
				continue
			}
			if ready[k] >= 2*f+1 {
				<-mm[index-1]
				outputChannel <- m.Data
				mm[index] <- true
				mu.Unlock()
				break
			}
		}
		mu.Unlock()
	}
	wg.Done()
}*/
