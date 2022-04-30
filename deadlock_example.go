package main

import (
	"fmt"
	"sync"
)

var w sync.WaitGroup

func inc(c chan int) {
	i := <-c
	i++
	c <- i
}

func sum(c chan int) int {
	sum := 0
	for i := range c {
		sum += i
	}
	return sum
}

func send(v int, c chan int) {
	c <- v
	w.Done()
}

func main() {
	c := make(chan int)
	w.Add(3)
	go send(5, c)
	go send(2, c)
	go send(17, c)
	w.Wait()
	close(c)
	fmt.Println(sum(c))

	/*
		w.Add(2)
		go inc(c)
		go inc(c)
		w.Wait()
		fmt.Println(<-c)
	*/
}
