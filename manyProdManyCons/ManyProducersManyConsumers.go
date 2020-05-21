package manyProdManyCons

import (
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var wgProducers sync.WaitGroup
var globalStopChannel chan bool

// Expected result: there are 2 producer goroutines that add numbers to a buffered queue channel.
// There are 2 consumer goroutines that consumes from the buffered queue channel.
// All goroutines (both producers and consumers) generate a random number each time they consume or produce.
// If that number is 0, they signal an orchestrator to kill all the producers.
// The consumers will die after all elements from the queue channel are consumed.
func ManyProducersManyConsumers() {
	wgProducers.Add(2)

	// Create communication channel, done channel and stop channel
	var queueChannel = make(chan int, 10)
	var doneChannel = make(chan bool)

	// Create global stop channel. The buffer size of 4 is the amount of producers + consumers. If we don't
	// assign a size bigger or equal to the producers + consumers, they will get blocked when trying to signal
	// the orchestrator.
	globalStopChannel = make(chan bool, 4)

	// Create 4 goroutines, 2 consumers and 2 producers.
	go producer(queueChannel, doneChannel)
	go producer(queueChannel, doneChannel)
	go consumer(queueChannel)
	go consumer(queueChannel)
	go orchestrator(doneChannel)

	// Wait until the reset is done
	wgProducers.Wait()

	// Close queueChannel so the producers can finish once they finish consuming all items.
	close(queueChannel)

	// This sleep is just to give time to the program to finish the goroutines. In a real environment the application
	// won't finish so we don't have the need of doing this.
	time.Sleep(4 * time.Second)
	println("Number of gorutines running: " + strconv.Itoa(runtime.NumGoroutine()))

}

func orchestrator(doneChannel chan<- bool) {
	select {
	case <-globalStopChannel:
		println("Stopping everything")
		close(doneChannel)
		return
	}
}

func producer(queueChannel chan int, doneChannel <-chan bool) {
	defer wgProducers.Done()

	for {
		// Select uses a pseudo random technique to choose between which case statement to check.
		// If we want to prioritize a statement, we can do it this way. We don't want to produce
		// more elements if we are told to stop.
		select {
		case <-doneChannel:
			println("Finishing producer")
			return
		default:
		}

		println("Adding 0 to queue")
		queueChannel <- 0
		getRandomNumberToStopEverything()
	}
}

func consumer(queueChannel <-chan int) {
	for value := range queueChannel {
		println("Consuming value from queue: " + strconv.Itoa(value))
		getRandomNumberToStopEverything()
	}

	println("Finishing consumer")

	return
}

func getRandomNumberToStopEverything() {
	// Returns a random number from 0 to 10. If 0, stop everything
	if randomNumber := rand.Intn(10); randomNumber == 0 {
		println("Someone randomed a 0. Signaling orchestrator to stop everything...")
		globalStopChannel <- true
	}
}
