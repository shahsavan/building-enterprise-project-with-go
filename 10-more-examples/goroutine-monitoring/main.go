package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof handlers
	"runtime"
	"time"
)

func main() {
	// A small worker pool (bounded concurrency).
	jobs := make(chan int, 1000)

	for i := 0; i < 8; i++ { // fixed number of workers
		go func(workerID int) {
			for job := range jobs {
				_ = job
				time.Sleep(5 * time.Millisecond) // simulate work
			}
		}(i)
	}

	// Generates work continuously.
	go func() {
		i := 0
		t := time.NewTicker(1 * time.Millisecond)
		defer t.Stop()

		for range t.C {
			i++

			// Normal mode: send to workers (stable goroutine count).
			select {
			case jobs <- i:
			default:
				// Queue full: backpressure; drop or block in real systems.
			}
		}
	}()

	// Optional: create goroutine buildup intentionally (to observe in pprof).
	// Visit: /start?leak=1
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("leak") == "1" {
			go func() {
				for {
					// This pattern is intentionally bad: one goroutine per task,
					// and the tasks block for a while.
					go func() {
						time.Sleep(30 * time.Second)
					}()
					time.Sleep(1 * time.Millisecond)
				}
			}()
			fmt.Fprintln(w, "started leak mode (intentional goroutine buildup)")
			return
		}
		fmt.Fprintln(w, "running normally (worker pool mode)")
	})

	// Simple health signal: goroutine count.
	http.HandleFunc("/health/goroutines", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "goroutines=%d\n", runtime.NumGoroutine())
	})

	// A tiny index page.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Try:")
		fmt.Fprintln(w, "  /health/goroutines")
		fmt.Fprintln(w, "  /debug/pprof/")
		fmt.Fprintln(w, "  /debug/pprof/goroutine?debug=1")
		fmt.Fprintln(w, "  /start            (normal mode)")
		fmt.Fprintln(w, "  /start?leak=1      (intentional goroutine buildup)")
	})

	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Printf("pprof at http://localhost%s/debug/pprof/", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
