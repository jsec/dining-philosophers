package main

import (
	"sync"
	"time"
)

func Philosopher(wg *sync.WaitGroup, id int, left, right *Chopstick, hostChan chan<- EatRequest) {
	defer wg.Done()

	for i := 0; i < MealMaximum; i++ {
		left.Lock()
		right.Lock()

		// Request permission to eat from the host
		hostResponse := make(chan func())
		hostChan <- EatRequest{
			id:       id,
			callback: hostResponse,
		}

		// Once we've received the appropriate callback response from the host, we can eat
		done := <-hostResponse
		logger(id, "is starting to eat")

		// Simulate a second for the time it takes to "eat"
		time.Sleep(1 * time.Second)

		logger(id, "is finished eating")
		right.Unlock()
		left.Unlock()
		done()
	}

	logger(id, "has finished all 3 courses")
}
