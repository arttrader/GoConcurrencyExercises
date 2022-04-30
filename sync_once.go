package main

import (
	"fmt"
	"sync"
)

var on sync.Once

func setup() {
	fmt.Println("init")
}

func dostuff() {
	// this needs to be completed first and only once
	on.Do(setup)
	fmt.Println("hello")
	wg.Done()
}

func main() {
	wg.Add(2)
	go dostuff()
	go dostuff()
	wg.Wait()
}
