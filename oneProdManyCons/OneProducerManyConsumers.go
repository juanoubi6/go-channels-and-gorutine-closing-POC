package oneProdManyCons

import (
	"runtime"
	"strconv"
	"sync"
	"time"
)

var producerWg sync.WaitGroup

// Expected result: there are 2 consumer goroutines that consume numbers from a buffered queue channel.
// There is one producer gouroutine that adds numbers to the buffered queue channel.
// When the producer goroutine produces 100 elements for the channel, it finishes and closes the queue channel.
// The consumers will consume from the queue until it's empty and then they will close.
func OneProducersManyConsumers() {
	producerWg.Add(2)

	// Create communication channel and done channel
	var queueChannel = make(chan int, 25)

	// Create 3 goroutines, one producer and 2 consumers.
	go producer(queueChannel)
	go consumer(queueChannel)
	go consumer(queueChannel)

	// Wait until the reset is done
	producerWg.Wait()

	println("Number of gorutines running: " + strconv.Itoa(runtime.NumGoroutine()))
}

// Produce 25 numbers and then close the queueChannel
func producer(queueChannel chan int) {
	for i := 0; i <= 25; i++ {
		println("Adding number to queue: " + strconv.Itoa(i))
		queueChannel <- i
	}

	// If I don't close the queue channel, the gorutines will keep waiting forever and there would be a deadlock.
	println("Finished producing numbers. Closing...")
	close(queueChannel)

	return
}

func consumer(queueChannel chan int) {
	defer func() {
		println("Closing consumer")
		producerWg.Done()
	}()

	for value := range queueChannel {
		println("Consuming value from queue: " + strconv.Itoa(value))
		time.Sleep(time.Second)
	}

	return
}
