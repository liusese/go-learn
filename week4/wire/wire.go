//+build wireinject

package main

import "github.com/google/wire"

func initEvent() Event {
	wire.Build(wire.NewSet(NewEvent, NewGreeter, NewMessage))
	return Event{}
}
