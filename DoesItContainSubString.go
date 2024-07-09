package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/schollz/progressbar/v3"
)

const chunkSize = 1024 * 1024 * 5 // 5MB, experiment with this size for optimal performance

func processChunk(filePath string, chunkStart, chunkSize int64, substring string, results chan<- string, wg *sync.WaitGroup, progress *int64, bar *progressbar.ProgressBar) {
	defer wg.Done()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	file.Seek(chunkStart, 0)
	reader := bufio.NewReader(file)

	buffer := make([]byte, chunkSize)
	n, err := reader.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		fmt.Printf("Failed to read file: %v\n", err)
		return
	}

	lines := strings.Split(string(buffer[:n]), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.Contains(line, substring) {
			results <- line
		}
	}

	atomic.AddInt64(progress, chunkSize)
	bar.Add64(chunkSize)
}

func chunkify(filePath string, chunkSize int64) ([][2]int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var chunks [][2]int64
	var chunkStart int64
	for chunkStart < fileInfo.Size() {
		chunkEnd := chunkStart + chunkSize
		if chunkEnd > fileInfo.Size() {
			chunkEnd = fileInfo.Size()
		}

		chunks = append(chunks, [2]int64{chunkStart, chunkEnd - chunkStart})
		chunkStart = chunkEnd
	}

	return chunks, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run checkpword.go <substring> <passwordlist>")
		return
	}

	startTime := time.Now()

	substring := os.Args[1]
	passwordListFile := os.Args[2]

	progress := int64(0)
	var wg sync.WaitGroup

	chunks, err := chunkify(passwordListFile, chunkSize)
	if err != nil {
		fmt.Printf("Failed to chunkify file: %v\n", err)
		return
	}

	fileSize, err := os.Stat(passwordListFile)
	if err != nil {
		fmt.Printf("Failed to get file size: %v\n", err)
		return
	}

	bar := progressbar.NewOptions64(
		fileSize.Size(),
		progressbar.OptionSetDescription("Loading passwords"),
		progressbar.OptionSetTheme(progressbar.Theme{Saucer: "=", SaucerHead: ">", BarStart: "[", BarEnd: "]"}),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionFullWidth(),
	)

	numWorkers := runtime.NumCPU() // Adjust the number of workers based on CPU cores
	sem := make(chan struct{}, numWorkers)
	results := make(chan string, 1000) // Buffered channel to collect results

	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	for _, chunk := range chunks {
		wg.Add(1)
		sem <- struct{}{}
		go func(chunkStart, chunkSize int64) {
			defer func() { <-sem }()
			processChunk(passwordListFile, chunkStart, chunkSize, substring, results, &wg, &progress, bar)
		}(chunk[0], chunk[1])
	}

	wg.Wait()
	close(results)
	bar.Finish()

	duration := time.Since(startTime)
	fmt.Printf("Total time taken: %v\n", duration)
}

