package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	yes := 0
	no := 0
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	for i := 0; i < 100; i++ {
		go func() {
			vote := requestVote()
			mu.Lock()
			defer mu.Unlock()
			if vote {
				yes++
			} else {
				no++
			}
			cond.Broadcast()
		}()
	}
	mu.Lock()
	for yes < 50 && no < 50 {
		cond.Wait()
	}
	if yes >= 50 {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
}

func requestVote() bool {
	return rand.Int()%2 == 0
}
