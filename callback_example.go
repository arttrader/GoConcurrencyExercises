package main

/*
	Example given in GopherCon 2018
*/

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

type Event struct{}
type Item struct {
	kind string
}

var items = map[string]Item{"a": Item{"item1"}, "b": Item{"item2"}}

func doSlowThing() { time.Sleep(10 * time.Millisecond) }

func Fetch(t string, f func(Item)) {
	go func() {
		item := items[t]
		doSlowThing()
		f(item)
	}()
}

func main() {
	start := time.Now()

	n := int32(0)
	Fetch("a", func(i Item) {
		fmt.Println(i)
		if atomic.AddInt32(&n, 1) == 2 {
			fmt.Println(time.Since(start))
			os.Exit(0)
		}
	})
	Fetch("b", func(i Item) {
		fmt.Println(i)
		if atomic.AddInt32(&n, 1) == 2 {
			fmt.Println(time.Since(start))
			os.Exit(0)
		}
	})

	select {}
}
