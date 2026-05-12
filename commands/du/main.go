package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type DuOptions struct {
	Human     bool
	Total     bool
	MaxDepth  int
	BlockSize int64
}

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		human     = flag.Bool("human", false, "человеко-читаемый вывод")
		total     = flag.Bool("c", false, "показать общий итог")
		maxDepth  = flag.Int("d", -1, "максимальная глубина")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("du - оценка использования дискового пространства")
		fmt.Println("Использование: du [-human] [-c] [-d глубина] [путь...]")
		fmt.Println("  -human    человеко-читаемый формат")
		fmt.Println("  -c        показать общий итог")
		fmt.Println("  -d        максимальная глубина рекурсии")
		return
	}

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	opts := DuOptions{
		Human:    *human,
		Total:    *total,
		MaxDepth: *maxDepth,
	}

	var totalSize int64
	for _, path := range paths {
		size := duPath(path, opts, 0)
		totalSize += size
	}

	if opts.Total && len(paths) > 1 {
		if opts.Human {
			fmt.Printf("%s\ttotal\n", formatBytes(uint64(totalSize)))
		} else {
			fmt.Printf("%d\ttotal\n", totalSize)
		}
	}
}

// duPath вычисляет размер директории
func duPath(path string, opts DuOptions, depth int) int64 {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return 0
	}

	if !info.IsDir() {
		size := info.Size()
		printSize(path, size, opts)
		return size
	}

	var total int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return 0
	}

	for _, entry := range entries {
		subPath := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if entry.IsDir() && (opts.MaxDepth == -1 || depth < opts.MaxDepth) {
			wg.Add(1)
			go func(p string, d int) {
				defer wg.Done()
				size := duPath(p, opts, d+1)
				mu.Lock()
				total += size
				mu.Unlock()
			}(subPath, depth)
		} else if !entry.IsDir() {
			size := info.Size()
			total += size
			if depth == 0 {
				printSize(subPath, size, opts)
			}
		}
	}

	wg.Wait()
	
	if depth == 0 || (opts.MaxDepth != -1 && depth <= opts.MaxDepth) {
		printSize(path, total, opts)
	}
	
	return total
}

// printSize выводит размер с путём
func printSize(path string, size int64, opts DuOptions) {
	if opts.Human {
		fmt.Printf("%s\t%s\n", formatBytes(uint64(size)), path)
	} else {
		fmt.Printf("%d\t%s\n", size, path)
	}
}

// formatBytes форматирует байты
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"K", "M", "G", "T", "P", "E"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}