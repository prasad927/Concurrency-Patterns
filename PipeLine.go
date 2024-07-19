package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)

	// go routine will stream the data into stream channel
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)

	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case val := <-stream:
				taken <- val
			}
		}
	}()
	return taken
}

func primeFinder(done <-chan int, randIntStream <-chan int) <-chan int {

	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}

	primes := make(chan int)

	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randomInt := <-randIntStream:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()

	return primes
}

func fanInPrimeFinderResult[T any](done <-chan int, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	fannedInStream := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}

	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	go func() {
		wg.Wait()
		close(fannedInStream)
	}()

	return fannedInStream
}

func PipeLineMain() {
	start := time.Now()
	done := make(chan int)
	defer close(done)
	randNumFetcher := func() int {
		return rand.Intn(500000000)
	}

	randIntStream := repeatFunc(done, randNumFetcher)

	// primeStream := primeFinder(done, randIntStream)

	// for rando := range take(done, primeStream, 10) {
	// 	fmt.Println(rando)
	// }

	// fan out
	cpuCount := runtime.NumCPU()
	primeFinderChannels := make([]<-chan int, cpuCount)
	for i := 0; i < cpuCount; i++ {
		primeFinderChannels[i] = primeFinder(done, randIntStream)
	}

	// fan in
	fannedInStream := fanInPrimeFinderResult(done, primeFinderChannels...)
	for num := range take(done, fannedInStream, 10) {
		fmt.Println(num)
	}
	fmt.Println(time.Since(start))
}