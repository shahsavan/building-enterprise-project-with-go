package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	printVehicleIDs([]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"})
}

func printVehicleIDs(vehicleIDs []string) {
	fmt.Println("goroutines before run:", runtime.NumGoroutine())
	for _, id := range vehicleIDs {
		go fmt.Println(id) // one goroutine per vehicle
	}
	fmt.Println("goroutines after run:", runtime.NumGoroutine())
	time.Sleep(time.Second) // fragile and unreliable
}
