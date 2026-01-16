package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type RideRequest struct {
	RiderID string
	Pickup  string
	Dropoff string
}

type Assignment struct {
	RequestID string
	VehicleID string
	Route     string
}

// StageSelectVehicle simulates matching ride requests to available vehicles.
func StageSelectVehicle(ctx context.Context, in <-chan RideRequest) <-chan Assignment {
	out := make(chan Assignment)
	go func() {
		defer close(out)
		for req := range in {
			asg := Assignment{
				RequestID: req.RiderID,
				VehicleID: fmt.Sprintf("vehicle-%s", req.RiderID),
				Route:     fmt.Sprintf("%s -> %s", req.Pickup, req.Dropoff),
			}

			select {
			case out <- asg:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// StageFormatMessage renders an assignment as JSON ready for publishing/logging.
func StageFormatMessage(ctx context.Context, in <-chan Assignment) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for asg := range in {
			payload, _ := json.Marshal(asg) // ignore error in this small demo

			select {
			case out <- string(payload):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	in := make(chan RideRequest)
	go func() {
		defer close(in)
		requests := []RideRequest{
			{RiderID: "alpha", Pickup: "airport", Dropoff: "hotel"},
			{RiderID: "bravo", Pickup: "downtown", Dropoff: "station"},
			{RiderID: "charlie", Pickup: "harbor", Dropoff: "conference"},
		}
		for _, req := range requests {
			in <- req
		}
	}()

	assignmentStream := StageSelectVehicle(ctx, in)
	messageStream := StageFormatMessage(ctx, assignmentStream)

	for msg := range messageStream {
		fmt.Println("assignment ready:", msg)
	}
}
