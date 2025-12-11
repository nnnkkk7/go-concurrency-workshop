// loggen generates JSON access log files for the concurrency workshop.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"
)

// LogEntry represents a single JSON access log entry.
type LogEntry struct {
	Timestamp      string `json:"timestamp"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	Status         int    `json:"status"`
	ResponseTimeMs int    `json:"response_time_ms"`
	Bytes          int    `json:"bytes"`
	UserID         string `json:"user_id"`
	IP             string `json:"ip"`
}

// Configuration for log generation
type Config struct {
	OutputDir    string
	FileCount    int
	LinesPerFile int
	Seed         uint64
	Verbose      bool
}

var (
	methods       = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	methodWeights = []int{70, 20, 5, 3, 2} // cumulative: GET=70%, POST=20%, etc.

	paths = []string{
		"/api/users",
		"/api/users/%d",
		"/api/products",
		"/api/products/%d",
		"/api/orders",
		"/api/orders/%d",
		"/api/auth/login",
		"/api/auth/logout",
		"/api/search",
		"/api/health",
		"/api/metrics",
		"/",
		"/static/%s",
	}
	pathWeights = []int{15, 10, 15, 10, 10, 5, 8, 3, 7, 5, 2, 5, 5}

	statuses      = []int{200, 201, 204, 301, 302, 400, 401, 403, 404, 500, 502, 503}
	statusWeights = []int{75, 5, 2, 1, 1, 3, 2, 1, 5, 3, 1, 1}
)

func main() {
	cfg := parseFlags()

	if err := run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.OutputDir, "output", "./logs", "Output directory for log files")
	flag.IntVar(&cfg.FileCount, "files", 200, "Number of log files to generate")
	flag.IntVar(&cfg.LinesPerFile, "lines", 50000, "Lines per file")
	flag.Uint64Var(&cfg.Seed, "seed", uint64(time.Now().UnixNano()), "Random seed for reproducibility")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Show progress during generation")
	flag.Parse()
	return cfg
}

func run(cfg *Config) error {
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Clean up existing log files before generating new ones
	pattern := filepath.Join(cfg.OutputDir, "*.json")
	existingFiles, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to search for existing files: %w", err)
	}
	for _, file := range existingFiles {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("failed to remove existing file %s: %w", file, err)
		}
	}
	if len(existingFiles) > 0 {
		fmt.Printf("Cleaned up %d existing log file(s)\n", len(existingFiles))
	}

	outputRoot, err := os.OpenRoot(cfg.OutputDir)
	if err != nil {
		return fmt.Errorf("failed to open output directory: %w", err)
	}
	defer outputRoot.Close()

	rng := rand.New(rand.NewPCG(cfg.Seed, cfg.Seed))

	fmt.Println("Generating log files...")
	startTime := time.Now()
	totalSize := int64(0)

	for i := 1; i <= cfg.FileCount; i++ {
		filename := fmt.Sprintf("access_%03d.json", i)

		size, err := generateLogFile(outputRoot, filename, cfg.LinesPerFile, rng)
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", filename, err)
		}

		totalSize += size

		if cfg.Verbose {
			fmt.Printf("  [%d/%d] %s (%d lines, %.1fMB)\n",
				i, cfg.FileCount, filename, cfg.LinesPerFile, float64(size)/(1024*1024))
		}
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nDone! Generated %d files (%.1fMB total) in %s\n",
		cfg.FileCount, float64(totalSize)/(1024*1024), elapsed.Round(time.Millisecond))

	return nil
}

func generateLogFile(root *os.Root, filename string, lineCount int, rng *rand.Rand) (int64, error) {
	file, err := root.Create(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)

	for i := 0; i < lineCount; i++ {
		entry := generateLogEntry(rng)
		if err := encoder.Encode(entry); err != nil {
			return 0, err
		}
	}

	info, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

func generateLogEntry(rng *rand.Rand) LogEntry {
	status := weightedRandom(statuses, statusWeights, rng)

	return LogEntry{
		Timestamp:      generateTimestamp(rng),
		Method:         weightedRandom(methods, methodWeights, rng),
		Path:           generatePath(rng),
		Status:         status,
		ResponseTimeMs: generateResponseTime(rng),
		Bytes:          generateBytes(status, rng),
		UserID:         fmt.Sprintf("user_%d", rng.IntN(1000000)),
		IP:             generateIP(rng),
	}
}

func generateTimestamp(rng *rand.Rand) string {
	// Range: 2025-01-10 00:00:00 to 2025-01-15 23:59:59 UTC
	start := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 15, 23, 59, 59, 0, time.UTC)

	delta := end.Unix() - start.Unix()
	sec := rng.Int64N(delta)
	t := start.Add(time.Duration(sec) * time.Second)

	ms := rng.IntN(1000)
	t = t.Add(time.Duration(ms) * time.Millisecond)

	return t.Format(time.RFC3339Nano)
}

func generatePath(rng *rand.Rand) string {
	pattern := weightedRandom(paths, pathWeights, rng)

	switch pattern {
	case "/api/users/%d":
		return fmt.Sprintf(pattern, rng.IntN(1000)+1)
	case "/api/products/%d":
		return fmt.Sprintf(pattern, rng.IntN(5000)+1)
	case "/api/orders/%d":
		return fmt.Sprintf(pattern, rng.IntN(10000)+1)
	case "/static/%s":
		files := []string{"app.js", "style.css", "logo.png", "favicon.ico", "bundle.js"}
		return fmt.Sprintf(pattern, files[rng.IntN(len(files))])
	default:
		return pattern
	}
}

func generateResponseTime(rng *rand.Rand) int {
	// Normal distribution: mean=100ms, stdDev=200ms
	mean := 100.0
	stdDev := 200.0

	// Box-Muller transform for normal distribution
	u1 := rng.Float64()
	u2 := rng.Float64()
	z := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)

	value := mean + stdDev*z

	// Clamp to 1-5000ms
	if value < 1 {
		value = 1
	}
	if value > 5000 {
		value = 5000
	}

	return int(value)
}

func generateBytes(status int, rng *rand.Rand) int {
	var min, max int

	switch {
	case status == 200:
		min, max = 100, 50000
	case status == 201:
		min, max = 50, 1000
	case status == 204:
		return 0
	case status == 301 || status == 302:
		min, max = 0, 100
	case status >= 400 && status < 500:
		min, max = 50, 500
	case status >= 500:
		min, max = 100, 1000
	default:
		min, max = 100, 10000
	}

	if min == max {
		return min
	}
	return min + rng.IntN(max-min)
}

func generateIP(rng *rand.Rand) string {
	// Generate private IP addresses
	choice := rng.IntN(3)

	switch choice {
	case 0:
		// 10.0.0.0/8
		return fmt.Sprintf("10.%d.%d.%d", rng.IntN(256), rng.IntN(256), rng.IntN(256))
	case 1:
		// 172.16.0.0/12
		return fmt.Sprintf("172.%d.%d.%d", 16+rng.IntN(16), rng.IntN(256), rng.IntN(256))
	default:
		// 192.168.0.0/16
		return fmt.Sprintf("192.168.%d.%d", rng.IntN(256), rng.IntN(256))
	}
}

func weightedRandom[T any](items []T, weights []int, rng *rand.Rand) T {
	// Calculate cumulative weights
	total := 0
	for _, w := range weights {
		total += w
	}

	// Random selection
	r := rng.IntN(total)
	cumulative := 0

	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return items[i]
		}
	}

	// Fallback (should never reach here)
	return items[len(items)-1]
}
