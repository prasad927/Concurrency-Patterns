package main

import (
	"fmt"
	"sync"
	"time"
)

/*
	Handling critical section without Mutex Lock
*/

func process(data int) int {
	time.Sleep(time.Second * 2)
	return data * 2
}

func processData(wg *sync.WaitGroup, result *int, data int) {
	defer wg.Done()

	processData := process(data)
	*result = processData

}

func ConfinementMain() {
	start := time.Now()
	var wg sync.WaitGroup

	input := []int{1, 2, 3, 4, 5, 6}
	result := make([]int, len(input))

	for i, data := range input {
		wg.Add(1)
		go processData(&wg, &result[i], data)
	}
	wg.Wait()

	fmt.Println(result)
	fmt.Println(time.Since(start))
}
