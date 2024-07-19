package main

import (
	"fmt"
	"sync"
)

func Producer(id int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < 3; i++ {
			ch <- id*10 + i // sending result in channel
		}
	}()

	return ch
}

func fanIn(inputs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	output := make(chan int)

	// Function to copy values from a channel to the output channel
	copy := func(c <-chan int) {
		defer wg.Done()
		for val := range c {
			output <- val
		}
	}

	wg.Add(len(inputs))
	for _, ch := range inputs {
		go copy(ch)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func FanInmain() {
	// Step-1 Producer
	ch1 := Producer(1)
	ch2 := Producer(20)

	// Step-1 FanIn
	mergedOutputChannel := fanIn(ch1, ch2) // combining results of multiple channels
	
	// Step-3 Consumer
	for val := range mergedOutputChannel {
		fmt.Println(val)
	}
}
