package observer

import (
	"context"
	"fmt"
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
