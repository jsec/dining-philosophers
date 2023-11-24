package main

import (
	"fmt"
	"sync"
)

const (
	PhilosopherCount    int = 5
	MealMaximum         int = 3
	MaxConcurrentEaters int = 2
)

type Chopstick struct{ sync.Mutex }

type EatRequest struct {
	id       int
	callback chan func()
}

type EatResponse struct {
	id   int
	done chan struct{}
}

func logger(id int, text string) {
	fmt.Printf("Philosopher %d %s\n", id, text)
}

func main() {
	chopsticks := make([]*Chopstick, PhilosopherCount)
	for i := 0; i < PhilosopherCount; i++ {
		chopsticks[i] = &Chopstick{}
	}

	// Create the main channel for philosophers to request permission to eat from the host
	hostChan := make(chan EatRequest)

	var wg sync.WaitGroup
	wg.Add(PhilosopherCount)

	// Spawn philosophers
	for i := 0; i < PhilosopherCount; i++ {
		go Philosopher(&wg, i, chopsticks[i], chopsticks[(i+1)%PhilosopherCount], hostChan)
	}

	// Spawn a separate goroutine to wait for the philosophers to finish
	go func() {
		wg.Wait()
		close(hostChan)
	}()

	Host(hostChan)
}
