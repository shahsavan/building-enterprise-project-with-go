package main

import "fmt"

func main() {
	deadlockExample()
	rangeDeadlock([]int{1, 2, 3})
	fixedExample()
}

func deadlockExample() {
	ch := make(chan int)
	ch <- 42 // fatal error: all goroutines are asleep - deadlock!
}

func rangeDeadlock(tasks []int) {
	ch := make(chan int)

	go func() {
		for _, t := range tasks {
			ch <- t
		}
		// channel is never closed
	}()

	for v := range ch {
		fmt.Println(v)
	}
	// blocks forever waiting for more values
}

func fixedExample() {
	// Buffered channel
	ch := make(chan int, 1)
	ch <- 42

	// Dedicated sender goroutine
	ch2 := make(chan int)
	go func() {
		ch2 <- 42
	}()
	value := <-ch2
	fmt.Println(value)
}
