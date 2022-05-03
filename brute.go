package main

import (
	"context"
	"flag"
	"fmt"
	"mikrobrute/workers"
	"sync"
)

func main() {
	numWorkers := flag.Int("w", 1, "number of workers")
	flag.Parse()

	jobChan := make(chan string, *numWorkers)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	fmt.Println(*numWorkers)
	for i := 1; i <= *numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			workers.New(i).ListenAndExecute(jobChan, ctx)
			cancel()
			wg.Done()
		}(i)
	}

	go func() {
		for i := 4000; i < 10000; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Received Cancellation")
				close(jobChan)
				return
			default:
				jobChan <- fmt.Sprintf("%04d", i)
			}
		}
		close(jobChan)
	}()
	fmt.Println("Waiting")
	wg.Wait()
	fmt.Println("Done")
}
