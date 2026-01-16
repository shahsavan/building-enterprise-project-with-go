package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type RideJob struct {
	ID      int
	Rider   string
	Pickup  string
	Dropoff string
}

type JobResult struct {
	JobID   int
	Vehicle string
	Err     error
}

func assignRide(job RideJob) JobResult {
	time.Sleep(100 * time.Millisecond) // simulate routing/dispatch time

	if job.Pickup == "maintenance-yard" {
		return JobResult{JobID: job.ID, Err: errors.New("no vehicles available at yard")}
	}

	return JobResult{
		JobID:   job.ID,
		Vehicle: fmt.Sprintf("vehicle-%02d", job.ID%4+1),
	}
}

func worker(ctx context.Context, id int, jobs <-chan RideJob, results chan<- JobResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return // channel closed, worker stops
			}

			fmt.Printf("worker %d dispatching rider %s\n", id, job.Rider)
			results <- assignRide(job)

		case <-ctx.Done():
			return // graceful stop when context is canceled
		}
	}
}

func StartPool(ctx context.Context, workerCount int) (chan<- RideJob, <-chan JobResult) {
	jobs := make(chan RideJob)
	results := make(chan JobResult, workerCount)

	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go worker(ctx, i, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return jobs, results
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	jobs, results := StartPool(ctx, 3)

	go func() {
		defer close(jobs)
		requests := []RideJob{
			{ID: 1, Rider: "alice", Pickup: "airport", Dropoff: "hotel"},
			{ID: 2, Rider: "bob", Pickup: "maintenance-yard", Dropoff: "garage"},
			{ID: 3, Rider: "chris", Pickup: "downtown", Dropoff: "stadium"},
			{ID: 4, Rider: "dana", Pickup: "harbor", Dropoff: "museum"},
			{ID: 5, Rider: "eric", Pickup: "suburb", Dropoff: "university"},
		}
		for _, job := range requests {
			select {
			case jobs <- job:
			case <-ctx.Done():
				return
			}
		}
	}()

	for res := range results {
		if res.Err != nil {
			fmt.Printf("job %d failed: %v\n", res.JobID, res.Err)
			continue
		}
		fmt.Printf("job %d assigned to %s\n", res.JobID, res.Vehicle)
	}
}
