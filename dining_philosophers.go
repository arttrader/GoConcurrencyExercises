package main

import (
	"fmt"
	"sync"
)

/*
Dining Philosophers Problem
• Classic problem involving concurrency and synchronization
Problem
• 5 philosophers sitting at a round table
• 1 chopstick is placed between each adjacent pair
• Want to eat rice from their plate, but needs two chopsticks
• Only one philosopher can hold a chopstick at a time
• Not enough chopsticks for everyone to eat at once
*/

var w sync.WaitGroup

type ChopS struct{ sync.Mutex }

type Philo struct {
	num             int
	leftCS, rightCS *ChopS
}

func (p Philo) eat() {
	for i := 0; i < 1; i++ {
		p.leftCS.Lock()
		p.rightCS.Lock()

		fmt.Printf("%d eating\n", p.num)

		p.rightCS.Unlock()
		p.leftCS.Unlock()
	}
	w.Done()
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
	/*
		creates deadlock if
		• All philosophers might lock their left chopsticks concurrently
		• None can lock their right chopsticks
	*/
	w.Add(5)
	for i := 0; i < 5; i++ {
		go philos[i].eat()
	}
	w.Wait()
}
