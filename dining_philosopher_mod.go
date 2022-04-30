package main

import (
	"fmt"
	"sync"
)

/*
Implement the dining philosopher’s problem with the following constraints/modifications.

There should be 5 philosophers sharing chopsticks, with one chopstick between each adjacent pair of philosophers.

Each philosopher should eat only 3 times (not in an infinite loop as we did in lecture)

The philosophers pick up the chopsticks in any order, not lowest-numbered first (which we did in lecture).

In order to eat, a philosopher must get permission from a host which executes in its own goroutine.

The host allows no more than 2 philosophers to eat concurrently.

Each philosopher is numbered, 1 through 5.

When a philosopher starts eating (after it has obtained necessary locks) it prints “starting to eat <number>” on a line by itself, where <number> is the number of the philosopher.

When a philosopher finishes eating (before it has released its locks) it prints “finishing eating <number>” on a line by itself, where <number> is the number of the philosopher.
*/

var w sync.WaitGroup

const ne int = 1

type ChopS struct{ sync.Mutex }

type Philo struct {
	num             int
	leftCS, rightCS *ChopS
}

type Host struct {
	nEating chan int
}

func (p Philo) eat(h *Host) {
	for i := 0; i < ne; i++ {
		if h.allow() {
			fmt.Printf("starting to eat %d \n", p.num)
			p.leftCS.Lock()
			p.rightCS.Lock()
			p.rightCS.Unlock()
			p.leftCS.Unlock()
			h.done()
			fmt.Printf("finishing eating %d \n", p.num)
		}
	}
	w.Done()
}

func (h Host) allow() bool {
	ne := <-h.nEating
	if ne < 2 {
		ne++
		h.nEating <- ne
		return true
	} else {
		return false
	}
}

func (h Host) done() {
	ne := <-h.nEating
	ne--
	h.nEating <- ne
}

func main() {
	CSticks := make([]*ChopS, 5)
	for i := 0; i < 5; i++ {
		CSticks[i] = new(ChopS)
	}
	var philos [5]*Philo
	for i := 0; i < 5; i++ {
		philos[i] = &Philo{i, CSticks[i], CSticks[(i+1)%5]}
	}
	var host *Host
	host = &Host{make(chan int)}
	host.nEating <- 0
	w.Add(5)
	for i := 0; i < 5; i++ {
		go philos[i].eat(host)
	}
	w.Wait()
}
