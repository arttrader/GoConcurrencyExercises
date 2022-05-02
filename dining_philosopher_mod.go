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

const nPhilos = 5
const ne = 3
const concurrentEatingNum = 2

var w sync.WaitGroup

type ChopS struct{ sync.Mutex }

type Philo struct {
	num             int
	leftCS, rightCS *ChopS
}

type Host struct {
	sync.Mutex
	nEating int
}

func (p Philo) eat(h *Host) {
	for i := 0; i < ne; i++ {
		if h.allow() {
			p.leftCS.Lock()
			p.rightCS.Lock()
			fmt.Printf("starting to eat %d \n", p.num)
			fmt.Printf("finishing eating %d \n", p.num)
			p.rightCS.Unlock()
			p.leftCS.Unlock()
			h.done()
		}
	}
	w.Done()
}

func (h Host) allow() bool {
	h.Lock()
	if h.nEating < concurrentEatingNum {
		h.nEating++
		h.Unlock()
		return true
	} else {
		h.Unlock()
		return false
	}
}

func (h Host) done() {
	h.Lock()
	h.nEating--
	h.Unlock()
}

func main() {
	CSticks := make([]*ChopS, nPhilos)
	for i := 0; i < nPhilos; i++ {
		CSticks[i] = new(ChopS)
	}
	var philos [nPhilos]*Philo
	for i := 0; i < nPhilos; i++ {
		philos[i] = &Philo{i, CSticks[i], CSticks[(i+1)%nPhilos]}
	}
	var host *Host
	host = new(Host)
	host.nEating = 0
	w.Add(nPhilos)
	for i := 0; i < nPhilos; i++ {
		go philos[i].eat(host)
	}
	w.Wait()
}
