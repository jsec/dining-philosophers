package main

import "slices"

func Host(hostChan <-chan EatRequest) {
	var currentlyEating []int
	receiveRequest := hostChan

	// Create channel to receive callbacks from philosophers when they've finished
	// eating
	callbackChan := make(chan EatResponse)

	// loop until we've successfully allowed a philosopher to eat
	for {
		select {
		case req, ok := <-receiveRequest:
			if !ok {
				return
			}

			currentlyEating = append(currentlyEating, req.id)

			req.callback <- func() {
				done := make(chan struct{})
				callbackChan <- EatResponse{id: req.id, done: done}
				// Wait to get a response (prevent greedy philosophers)
				<-done
			}
		case req := <-callbackChan:
			// remove the philosopher's id from the currently eating slice
			idx := slices.Index(currentlyEating, req.id)
			currentlyEating = append(currentlyEating[:idx], currentlyEating[idx+1:]...)

			// Close the ephemeral channel created when the host granted permission to eat
			close(req.done)
		}

		if len(currentlyEating) < MaxConcurrentEaters {
			receiveRequest = hostChan
		} else {
			// Setting the channel to nil ensures that all requests are ignored until
			// we process a message from the callback chan and free up room for another
			// philosopher to eat
			receiveRequest = nil
		}
	}
}
