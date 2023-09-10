package main

import (
	"fmt"
)

type Consumer struct {
	msgs chan int
}

func NewCustomer(msgs chan int) *Consumer {
	return &Consumer{msgs: msgs}
}

func (c *Consumer) consume() {
	fmt.Println("consume: Started")
	for {
		msg := <-c.msgs
		fmt.Println("consumer receive: ", msg)
	}
}

type Producer struct {
	msgs chan int
	done *chan bool
}

func NewProducer(msgs chan int, done *chan bool) *Producer {
	return &Producer{msgs: msgs, done: done}
}

func (p *Producer) produce(max int) {
	fmt.Println("produce: Started")
	for i := 0; i < max; i++ {
		fmt.Println("producer send: ", i)
		p.msgs <- i
	}
	*p.done <- true
	fmt.Println("produce: Done")
}

func main() {
	var msgs = make(chan int)
	var done = make(chan bool)
	go NewProducer(msgs, &done).produce(10)
	go NewCustomer(msgs).consume()
	<-done
}
