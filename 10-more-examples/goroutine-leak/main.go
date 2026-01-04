// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	leakExample()
}

func leakExample() {
	fmt.Println("goroutines before leak:", runtime.NumGoroutine())
	out := make(chan int)

	go func() {
		for i := 0; ; i++ {
			out <- i // blocks when no receiver exists
		}
	}()

	for n := range out {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
	time.Sleep(time.Second)
	fmt.Println("goroutines after leak:", runtime.NumGoroutine())
}
