package main

/*
	Example given in GopherCon 2018
*/

import (
	"fmt"
	"time"
)

type Event struct{}
type Item struct {
	kind string
}

var items = map[string]Item{"a": Item{"item1"}, "b": Item{"item2"}}

func doSlowThing() { time.Sleep(10 * time.Millisecond) }

func Fetch(t string) <-chan Item {
	c := make(chan Item, 1)
	go func() {
		item, ok := items[t]
		doSlowThing()
		if !ok {
			close(c)
			return
		}
		c <- item
	}()
	return c
}

func consume(a, b Item) {
	fmt.Println(a, b)
}

func main() {
	start := time.Now()
	a := Fetch("a")
	b := Fetch("b")
	consume(<-a, <-b)
	fmt.Println(time.Since(start))
}
