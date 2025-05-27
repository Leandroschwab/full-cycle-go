package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

func main() {
	// Parse command line arguments
	url := flag.String("url", "", "URL of the service to be tested")
	requests := flag.Int("requests", 0, "Number of total requests")
	concurrency := flag.Int("concurrency", 0, "Number of concurrent calls")
	flag.Parse()

	// Validate input parameters
	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("All parameters are required and must be valid")
		fmt.Println("Usage: --url=<url> --requests=<requests> --concurrency=<concurrency>")
		os.Exit(1)
	}

	// Execute the stress test
	fmt.Printf("Starting stress test for %s\n", *url)
	fmt.Printf("Total requests: %d\n", *requests)
	fmt.Printf("Concurrency level: %d\n", *concurrency)
	fmt.Println("--------------------------------------------------")

	startTime := time.Now()
	results := executeStressTest(*url, *requests, *concurrency)
	totalDuration := time.Since(startTime)

	// Generate and print report
	printReport(results, totalDuration, *requests)
}

func executeStressTest(url string, totalRequests, concurrency int) []Result {
	results := make([]Result, 0, totalRequests)
	jobs := make(chan int, totalRequests)
	resultsChan := make(chan Result, totalRequests)
	var wg sync.WaitGroup

	// Create worker goroutines
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(url, jobs, resultsChan, &wg)
	}

	// Send jobs to workers
	for i := 0; i < totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(resultsChan)

	// Collect results
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func worker(url string, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for range jobs {
		startTime := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(startTime)

		result := Result{
			Duration: duration,
			Error:    err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			resp.Body.Close()
		}

		results <- result
	}
}

func printReport(results []Result, totalDuration time.Duration, totalRequests int) {
	// Count status codes
	statusCodes := make(map[int]int)
	var successCount int
	var failedCount int

	for _, result := range results {
		if result.Error != nil {
			failedCount++
		} else {
			statusCodes[result.StatusCode]++
			if result.StatusCode == 200 {
				successCount++
			}
		}
	}

	// Print report
	fmt.Println("--------------------------------------------------")
	fmt.Println("Stress Test Report")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total time: %v\n", totalDuration)
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Successful requests (HTTP 200): %d\n", successCount)
	fmt.Println("Status code distribution:")
	for code, count := range statusCodes {
		fmt.Printf("  HTTP %d: %d\n", code, count)
	}
	if failedCount > 0 {
		fmt.Printf("Failed requests: %d\n", failedCount)
	}
	fmt.Println("--------------------------------------------------")
}
