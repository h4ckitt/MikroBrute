package main

import (
	"flag"
	"fmt"
	"mikrobrute/workers"
	"sync"
)

func main() {
	/*fmt.Println("vim-go")
	resp, err := http.Get("http://ictplaza.hs/login")

	if err != nil {
		log.Fatalln(err)
		return
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))*/

	numWorkers := flag.Int("w", 1, "number of workers")

	jobChan := make(chan string, *numWorkers)

	/*	manager := workers.NewPool(*numWorkers, jobChan)
		manager.Start()*/
	wg := sync.WaitGroup{}
	for i := 1; i <= *numWorkers; i++ {
		go func(i int) {
			workers.New(i).ListenAndExecute(jobChan)
			wg.Done()
		}(i)
	}

	for i := 4018; i < 10000; i++ {
		jobChan <- fmt.Sprintf("%04d", i)
	}
	close(jobChan)
	wg.Wait()
	fmt.Println("Done")
}
