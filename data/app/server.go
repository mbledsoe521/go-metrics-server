package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	port            = 12345
	metricsFilePath = "/data/metrics_from_special_app.txt"
	cacheExpiration = 10 * time.Second
)

var (
	mu          sync.Mutex
	metrics     map[string]string
	lastUpdated time.Time
)

func main() {
	// Start a goroutine to update the metrics cache periodically
	go func() {
		for {
			updateMetricsCache()
			time.Sleep(cacheExpiration)
		}
	}()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Acquire lock to ensure thread-safe access to metrics
		mu.Lock()
		defer mu.Unlock()

		// Check if cache is expired
		if time.Since(lastUpdated) > cacheExpiration {
			updateMetricsCache()
		}

		// Write metrics to response
		for key, value := range metrics {
			fmt.Fprintf(w, "%s=%s\n", key, value)
		}
	})

	fmt.Printf("Server listening on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func updateMetricsCache() {
	mu.Lock()
	defer mu.Unlock()

	// Open metrics file
	file, err := openFile(metricsFilePath)
	if err != nil {
		fmt.Printf("Error opening metrics file: %v\n", err)
		return
	}
	defer file.Close()

	// Read metrics line by line
	metrics = make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid metric line: %s\n", line)
			continue
		}
		metrics[parts[0]] = parts[1]
	}

	lastUpdated = time.Now()
}

func openFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create the file if it doesn't exist
			file, err = os.Create(path)
		}
	}
	return file, err
}
