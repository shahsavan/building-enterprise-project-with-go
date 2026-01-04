package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	runTimedTask()
}

func runTimedTask() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	go func() {
		<-ctx.Done()
		fmt.Println("background task stopped")
	}()

	// function exits without canceling the context
}

func loadRoute(routeID string) (*Route, error) {
	resp, err := http.Get(fmt.Sprintf(
		"https://routing.service/api/routes/%s", routeID,
	))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeRoute(resp.Body)
}

// Route represents a route object.
type Route struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// decodeRoute decodes the response body into a Route object.
func decodeRoute(body io.ReadCloser) (*Route, error) {
	var route Route
	if err := json.NewDecoder(body).Decode(&route); err != nil {
		return nil, err
	}
	return &route, nil
}
