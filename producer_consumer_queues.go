package main

/*
	Example given in GopherCon 2018
*/

import (
	"fmt"
	"path/filepath"
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
		item := items[t]
		c <- item
	}()
	return c
}

func Glob(pattern string) <-chan Item {
	c := make(chan Item)
	go func() {
		defer close(c)
		for t, item := range items {
			if ok, _ := filepath.Match(pattern, t); !ok {
				continue
			}
			c <- item
		}
	}()
	return c
}

func main() {
	for item := range Glob("[ab]*") {
		fmt.Println(item)
	}
}
