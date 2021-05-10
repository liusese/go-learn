package main

import "fmt"

type Message string

type Greeter struct {
	Message Message
}

type Event struct {
	Greeter Greeter
}

func NewMessage() Message {
	return Message("Hi there!")
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

func (g Greeter) Greet() Message {
	return g.Message
}

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}

func (e Event) Strat() {
	msg := e.Greeter.Greet()
	fmt.Println(msg)
}

func main() {
	// msg := NewMessage()
	// g := NewGreeter(msg)
	// e := NewEvent(g)
	e := initEvent()

	e.Strat()
	fmt.Println("over")
}
