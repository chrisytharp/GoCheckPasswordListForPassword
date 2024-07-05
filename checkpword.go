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

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/schollz/progressbar/v3"
)

const chunkSize = 1024 * 1024 * 50 // 50MB, experiment with this size for optimal performance

func processChunk(filePath string, chunkStart, chunkSize int64, bloomFilter *bloom.BloomFilter, wg *sync.WaitGroup, progress *int64, bar *progressbar.ProgressBar) {
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
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return
	}

	lines := strings.Split(string(buffer[:n]), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			bloomFilter.AddString(line)
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
		fmt.Println("Usage: go run checkpword.go <password> <passwordlist>")
		return
	}

	password := os.Args[1]
	passwordListFile := os.Args[2]

	bloomFilter := bloom.NewWithEstimates(10000000, 0.01)
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

	numWorkers := runtime.NumCPU() * 2 // Adjust the number of workers based on CPU cores
	sem := make(chan struct{}, numWorkers)

	for _, chunk := range chunks {
		wg.Add(1)
		sem <- struct{}{}
		go func(chunkStart, chunkSize int64) {
			defer func() { <-sem }()
			processChunk(passwordListFile, chunkStart, chunkSize, bloomFilter, &wg, &progress, bar)
		}(chunk[0], chunk[1])
	}

	wg.Wait()
	bar.Finish()

	fmt.Printf("Checking password: %s\n", password)
	if bloomFilter.TestString(password) {
		fmt.Println("Password might be in the list (with a low probability of false positives).")
	} else {
		fmt.Println("Password not in the list.")
	}
}
