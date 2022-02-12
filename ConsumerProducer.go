package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var cond sync.Cond

func producer(buffer chan<- int) {
	for {
		//Lock first
		cond.L.Lock()
		//checking if the buffer is full
		fmt.Printf("Checking if the buffer is full\n")
		fmt.Printf("Buffer size :%d \n", len(buffer))
		for len(buffer) == 5 {
			fmt.Println("Buffer is full, producer going to sleep")
			cond.Wait()
		}
		//add item in the buffer
		num := rand.Intn(800)
		fmt.Printf("Producer is awake, adding %d to the buffer \n", num)
		buffer <- num
		fmt.Printf("Buffer size :%d \n", len(buffer))
		// Access to the public area ends, and printing ends, unlock
		cond.L.Unlock()
		// Wake up consumers blocked on condition variables
		cond.Signal()
	}
}

func consumer(buffer <-chan int) {
	for {
		// Lock first
		cond.L.Lock()
		// Determine whether the buffer is empty
		fmt.Printf("Checking if the buffer is empty\n")
		fmt.Printf("Buffer size :%d\n", len(buffer))
		for len(buffer) == 0 {
			fmt.Println("Buffer is empty, consumer going to sleep")
			cond.Wait()
		}
		fmt.Println("Reading values from the buffer")
		num := <-buffer
		fmt.Printf("Consumer is awake and reading the value %d from the buffer.\n", num)
		// After accessing the public area, unlock
		cond.L.Unlock()
		// Wake up producers blocked on condition variables
		cond.Signal()
	}
}

func main() {
	ch := make(chan int, 5)
	cond.L = new(sync.Mutex)

	go producer(ch)
	for i := 0; i < 5; i++ {
		go consumer(ch)
	}
	time.Sleep(time.Second * 1)
}
