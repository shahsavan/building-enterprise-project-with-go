package main

import (
	"context"
	"fmt"
	"time"
)

type TransportEvent[T any] struct {
	Name string
	Data T
}

type TransportObservable[T any] struct {
	subs []chan TransportEvent[T]
}

func (o *TransportObservable[T]) Subscribe() <-chan TransportEvent[T] {
	ch := make(chan TransportEvent[T], 16) // small buffer for short bursts

	o.subs = append(o.subs, ch)

	return ch
}

func (o *TransportObservable[T]) Publish(evt TransportEvent[T]) {

	for _, ch := range o.subs {
		select {
		case ch <- evt:
		default:
			// observer is slow; choose drop, block, or log
			// depending on transport system requirements
		}
	}
}

func (o *TransportObservable[T]) Close() {
	for _, ch := range o.subs {
		close(ch)
	}
}

func StartObserver[T any](ctx context.Context, id int, events <-chan TransportEvent[T]) {
	go func() {
		for {
			select {
			case evt, ok := <-events:
				if !ok {
					return
				}
				fmt.Printf("observer %d received: %s\n", id, evt.Name)

			case <-ctx.Done():
				return
			}
		}
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	obs := &TransportObservable[string]{}

	// Register two observers
	events1 := obs.Subscribe()
	events2 := obs.Subscribe()

	StartObserver(ctx, 1, events1)
	StartObserver(ctx, 2, events2)

	// Publish a few events
	obs.Publish(TransportEvent[string]{Name: "created", Data: "order-123"})
	obs.Publish(TransportEvent[string]{Name: "processed", Data: "order-123"})
	obs.Publish(TransportEvent[string]{Name: "shipped", Data: "order-123"})

	// Give observers time to consume before shutting down
	time.Sleep(100 * time.Millisecond)
	obs.Close()
}
