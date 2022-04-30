package main

import (
	"fmt"
	"sync"
	"time"
)

var mut sync.Mutex

var i int

func calc1() {
	fmt.Println("calc 1")
	mut.Lock()
	i = i + 5
	mut.Unlock()
}

func calc2() {
	fmt.Println("calc 2")
	mut.Lock()
	i = i + 1
	mut.Unlock()
}

func print() {
	fmt.Printf("i: %d\n", i)
}

func main() {
	i = 0
	go calc1()
	go calc2()
	go print()

	fmt.Println("waiting...")
	time.Sleep(time.Second)
}
