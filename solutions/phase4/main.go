package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/nnnkkk7/go-concurrency-workshop/pkg/logparser"
)

func main() {
	startTime := time.Now()

	logDir := "../../logs"

	logRoot, err := os.OpenRoot(logDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log directory: %v\n", err)
		os.Exit(1)
	}
	defer logRoot.Close()

	pattern := filepath.Join(logDir, "access_*.json")
	fullPaths, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding log files: %v\n", err)
		os.Exit(1)
	}

	files := make([]string, len(fullPaths))
	for i, path := range fullPaths {
		files[i] = filepath.Base(path)
	}

	numWorkers := runtime.NumCPU()
	results := processFiles(logRoot, files, numWorkers)

	printResults(results, time.Since(startTime))
}

// processFiles は最適化されたワーカープールパターンでファイルを処理します
func processFiles(root *os.Root, files []string, numWorkers int) []*logparser.Result {
	fileCount := len(files)

	// ジョブ用にはより小さいバッファを使用（保留中の作業のみを保持）
	jobs := make(chan string, numWorkers*2)
	// 結果用にはより大きいバッファを使用（全ての結果を保持）
	results := make(chan *logparser.Result, fileCount)

	var wg sync.WaitGroup

	// ワーカーを起動
	for i := 0; i < numWorkers; i++ {
		wg.Go(func() {
			for filename := range jobs {
				result, err := processFile(root, filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", filename, err)
					continue
				}
				results <- result
			}
		})
	}

	// ジョブを供給（バッファ付きチャネルでノンブロッキング）
	go func() {
		for _, filename := range files {
			jobs <- filename
		}
		close(jobs)
	}()

	// 全てのワーカーが完了するのを待つ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 正確な容量で結果スライスを事前割り当て
	resultList := make([]*logparser.Result, 0, fileCount)
	for result := range results {
		resultList = append(resultList, result)
	}

	return resultList
}

// processFile は1つのログファイルを解析します
func processFile(root *os.Root, filename string) (*logparser.Result, error) {
	file, err := root.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &logparser.Result{
		FileName:     filename,
		StatusCounts: make(map[int]int),
	}

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var entry logparser.LogEntry
		if err := decoder.Decode(&entry); err != nil {
			continue
		}
		result.TotalCount++
		result.StatusCounts[entry.Status]++
	}

	return result, nil
}

// printResults は処理結果を表示します
func printResults(results []*logparser.Result, elapsed time.Duration) {
	totalRequests := 0
	totalStatusCounts := make(map[int]int)

	for _, result := range results {
		totalRequests += result.TotalCount
		for status, count := range result.StatusCounts {
			totalStatusCounts[status] += count
		}
	}

	fmt.Printf("\n=== 処理結果 ===\n")
	fmt.Printf("処理時間: %.2f秒\n", elapsed.Seconds())
	fmt.Printf("総リクエスト数: %s件\n", formatNumber(totalRequests))
	fmt.Printf("\nステータスコード別:\n")
	for status := 200; status <= 599; status += 100 {
		for s := status; s < status+100; s++ {
			if count, ok := totalStatusCounts[s]; ok {
				percentage := float64(count) / float64(totalRequests) * 100
				fmt.Printf("  %d: %s件 (%.2f%%)\n", s, formatNumber(count), percentage)
			}
		}
	}

	errorCount := 0
	for status, count := range totalStatusCounts {
		if status >= 400 {
			errorCount += count
		}
	}
	errorRate := float64(errorCount) / float64(totalRequests) * 100
	fmt.Printf("\nエラー率 (4xx, 5xx): %.2f%%\n", errorRate)
}

// formatNumber は数値を3桁カンマ区切りでフォーマットします
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}
