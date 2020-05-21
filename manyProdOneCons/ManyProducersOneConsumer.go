package manyProdOneCons

import (
	"runtime"
	"strconv"
	"sync"
	"time"
)

var reset sync.WaitGroup

// Expected result: there are 2 producer goroutines that add numbers to a buffered queue channel.
// There is a consumer goroutine that consumes from the buffered queue channel.
// When the consumer goroutine consumes 23 elements from the channel, it adds a value to the WaitGroup so the
// program can finish. The goroutines are closed before exiting.
func ManyProducersOneConsumer() {
	reset.Add(1)

	// Create communication channel and done channel
	var queueChannel = make(chan int, 10)
	var doneChannel = make(chan bool)

	// Create 3 goroutines, one consumer and 2 producers.
	go produceZero(queueChannel, doneChannel)
	go produceOne(queueChannel, doneChannel)
	go consumer(queueChannel, doneChannel)

	// Wait until the reset is done
	reset.Wait()

	// This sleep is just to give time to the program to finish the goroutines. In a real environment the application
	// won't finish so we don't have the need of doing this.
	time.Sleep(2 * time.Second)
	println("Number of gorutines running: " + strconv.Itoa(runtime.NumGoroutine()))

}

func produceZero(queueChannel chan int, doneChannel <-chan bool) {
	for {
		select {
		case <-doneChannel:
			println("Finishing 0s goroutine")
			return
		default:
			println("Adding 0 to queue")
			queueChannel <- 0
		}
	}
}

func produceOne(queueChannel chan int, doneChannel <-chan bool) {
	for {
		// Select uses a pseudo random technique to choose between which case statement to check.
		// If we want to prioritize a statement, we can do it this way. We don't want to produce
		// more elements if we are told to stop.
		select {
		case <-doneChannel:
			println("Finishing 1s goroutine")
			return
		default:
		}

		println("Adding 1 to queue")
		queueChannel <- 1
	}
}

//Consume from the queueChannel and after consuming 23 numbers, add a value to the WaitGroup.
func consumer(queueChannel chan int, doneChannel chan bool) {
	defer reset.Done()

	var counter = 0
	for value := range queueChannel {
		println("Consuming value from queue: " + strconv.Itoa(value))
		counter++
		if counter == 23 {
			println("Consumed 23 values from queue Channel, stop consuming.")
			close(doneChannel)
			return
		}
	}
}
