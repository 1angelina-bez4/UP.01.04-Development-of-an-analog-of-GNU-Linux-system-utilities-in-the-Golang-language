package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Counts struct {
	Lines int
	Words int
	Bytes int
}

func main() {
	// Определение флагов
	var (
		help   = flag.Bool("h", false, "показать справку")
		lines  = flag.Bool("l", false, "только строки")
		words  = flag.Bool("w", false, "только слова")
		bytes  = flag.Bool("c", false, "только байты")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("wc - подсчёт строк, слов и байт")
		fmt.Println("Использование: wc [-l] [-w] [-c] [файлы...]")
		fmt.Println("  -l    только строки")
		fmt.Println("  -w    только слова")
		fmt.Println("  -c    только байты")
		return
	}

	// Если не указаны флаги, показываем всё
	showAll := !(*lines || *words || *bytes)
	if !showAll {
		showAll = !(*lines || *words || *bytes)
	}

	files := flag.Args()
	var total Counts

	if len(files) == 0 {
		counts := processFile("-", os.Stdin)
		printCounts(counts, "-", showAll, *lines, *words, *bytes)
		total = addCounts(total, counts)
	} else {
		for _, fname := range files {
			file, err := os.Open(fname)
			if err != nil {
				fmt.Fprintf(os.Stderr, "wc: %v\n", err)
				continue
			}
			counts := processFile(fname, file)
			printCounts(counts, fname, showAll, *lines, *words, *bytes)
			total = addCounts(total, counts)
			file.Close()
		}
	}

	// Выводим общий итог
	if len(files) > 1 {
		printCounts(total, "total", showAll, *lines, *words, *bytes)
	}
}

// processFile подсчитывает строки, слова и байты в файле
func processFile(fname string, file *os.File) Counts {
	var counts Counts
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := scanner.Text()
		counts.Lines++
		counts.Bytes += len(line) + 1 // +1 для новой строки
		counts.Words += len(strings.Fields(line))
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "wc: ошибка чтения %s: %v\n", fname, err)
	}
	
	return counts
}

// printCounts выводит подсчитанные значения
func printCounts(counts Counts, name string, showAll, showLines, showWords, showBytes bool) {
	if showLines || showAll {
		fmt.Printf("%8d ", counts.Lines)
	}
	if showWords || showAll {
		fmt.Printf("%8d ", counts.Words)
	}
	if showBytes || showAll {
		fmt.Printf("%8d ", counts.Bytes)
	}
	fmt.Printf(" %s\n", name)
}

// addCounts суммирует два набора подсчётов
func addCounts(a, b Counts) Counts {
	return Counts{
		Lines: a.Lines + b.Lines,
		Words: a.Words + b.Words,
		Bytes: a.Bytes + b.Bytes,
	}
}