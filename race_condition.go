package main

import (
	"fmt"
	"time"
)

/*
When running goroutines concurrently and variables or anything is shared between them (communication between them), you can get very unpredictable results.

The example code demonstrates this by sharing the variable x between the two goroutines. Sometimes you get 0, sometimes you get 2 printed as output, depending on which goroutine gets completed first.
*/

func main() {
	var x int

	for i := 0; i < 5; i++ {
		x = 0
		go func() {
			fmt.Println("goroutine 1")
			x = 1
			x = x + 1
		}()
		go func() {
			fmt.Println("goroutine 2")
			fmt.Printf("x: %d\n", x)
		}()
		fmt.Println("waiting...")
		time.Sleep(time.Second)
	}
	fmt.Println("done")
}
