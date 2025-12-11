package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

	results := processFiles(logRoot, files)

	printResults(results, time.Since(startTime))
}

// ============================================================
// TODO: 以下の関数を自由に最適化してください
// ============================================================
// 目標: Phase 3よりもさらに高速化する
// 制約: なし（あらゆる最適化手法を試してください）
//
// 最適化のアイデア:
// - processFiles() の並行処理戦略を改善
// - processFile() のI/O処理を最適化
// - チャネルやgoroutineの使い方を工夫
// - メモリアロケーションを削減
// - バッファサイズを調整
// など、何でも試してみましょう！

func processFiles(root *os.Root, files []string) []*logparser.Result {
	// TODO: ここに実装を書いてください
	return nil
}

func processFile(root *os.Root, filename string) (*logparser.Result, error) {
	// TODO: 必要に応じてこの関数も最適化してください
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

// ============================================================
// 以下の関数は変更不要です
// ============================================================

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

// formatNumber は数値を3桁カンマ区切りでフォーマットします（変更不要）
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
