package main

import (
	"fmt"
	"sync"
	"time"
)


// place from some task will be going to generate
func source(tasks chan<- int, numTasks int) {
	for i := 0; i < numTasks; i++ {
		tasks <- i
	}
	close(tasks)
}

func worker(id int, taskQueue <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskQueue {
		// Simulating task processing
		fmt.Printf("Worker %d processing task %d\n", id, task)
		time.Sleep(time.Millisecond * time.Duration(task))
		results <- task * 2 // task 1 * 2 | task 2 * 2
	}
}

func fanOutMain() {
	const numWorkers = 3  // number of workers will be used to process tasks from source
	const numTasks = 10

	taskQueue := make(chan int, numTasks)
	resultQueue := make(chan int, numTasks)

	var wg sync.WaitGroup

	// start source go-routing to generate tasks
	go source(taskQueue, numTasks)

	// create worker go-routine
	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go worker(i, taskQueue, resultQueue, &wg)
	}

	//wait for all workers to finish
	go func() {
		wg.Wait()
		close(resultQueue)
	}()

	 // Process results from the result queue
	for result := range resultQueue {
		fmt.Printf("Result received: %d\n", result)
	}
	  
	fmt.Println("All tasks completed.")
}
