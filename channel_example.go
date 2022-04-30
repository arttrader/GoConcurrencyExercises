package main

import (
	"fmt"
	"sync"
)

var w sync.WaitGroup

func prod(v1 int, v2 int, c chan int) {
	c <- v1 * v2
}

func inc(c chan int) {
	i := <-c
	i++
	c <- i
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	go prod(1, 5, c1)
	go prod(3, 4, c2)
	// only wait for first data
	select {
	case a := <-c1:
		fmt.Println(a)
	case b := <-c2:
		fmt.Println(b)
	}

	// deadlock example
	c := make(chan int)
	w.Add(2)
	go inc(c)
	go inc(c)
	w.Wait()
	fmt.Println(<-c)
}
