package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// result represents the result of a URL check
type result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Error      error
}

func checkURL(client *http.Client, url string) result {
	start := time.Now()

	resp, err := client.Get(url)
	if err != nil {
		return result{URL: url, Duration: time.Since(start), Error: err}
	}
	defer resp.Body.Close()

	return result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   time.Since(start),
	}
}

// worker continuously receives url jobs from the jobs channel and sends the results to the results channel
func worker(id int, jobs <-chan string, results chan<- result, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		fmt.Printf("Worker %d checking URL: %s\n", id, url)
		results <- checkURL(client, url)
	}
}

func main() {
	urls := []string{
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://payd.com",
		"http://invalid.com",
	}
	const workerCount = 3

	// one shared http client for all workers. timeout prevents hanging requests
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// channels coordinate communication between workers and main goroutine
	jobs := make(chan string, len(urls))
	results := make(chan result, len(urls))

	var wg sync.WaitGroup

	// start workers
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i, jobs, results, client, &wg)
	}

	// send urls to jobs channel and close the jobs channel
	for _, url := range urls {
		jobs <- url
	}
	close(jobs)

	// wait for all workers to finish before closing the results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// receive results from the results channel and print them
	var successful, failed int

	fmt.Println("\nURL CHECK RESULTS")
	fmt.Println("------------------")

	for r := range results {
		if r.Error != nil {
			fmt.Printf("[ERROR] %s | %v\n", r.URL, r.Error)
			failed++
			continue
		}

		fmt.Printf("[%d] %s | %v\n", r.StatusCode, r.URL, r.Duration)
		successful++
	}

	fmt.Println("\nSUMMARY")
	fmt.Println("--------")
	fmt.Printf("Total:      %d\n", len(urls))
	fmt.Printf("Successful: %d\n", successful)
	fmt.Printf("Failed:     %d\n", failed)
}
