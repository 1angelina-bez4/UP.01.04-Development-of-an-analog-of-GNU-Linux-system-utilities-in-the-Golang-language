package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		number    = flag.Bool("n", false, "нумеровать строки")
		numberNonBlank = flag.Bool("b", false, "нумеровать только непустые строки")
		squeeze   = flag.Bool("s", false, "объединять несколько пустых строк в одну")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("cat - вывод содержимого файлов")
		fmt.Println("Использование: cat [-n] [-b] [-s] [файлы...]")
		fmt.Println("  -n    нумеровать все строки")
		fmt.Println("  -b    нумеровать только непустые строки")
		fmt.Println("  -s    объединять пустые строки")
		return
	}

	// Если файлы не указаны, читаем stdin
	if flag.NArg() == 0 {
		processFile("-", *number, *numberNonBlank, *squeeze)
		return
	}

	// Обрабатываем каждый файл
	for _, fname := range flag.Args() {
		if err := processFile(fname, *number, *numberNonBlank, *squeeze); err != nil {
			fmt.Fprintf(os.Stderr, "cat: %v\n", err)
			os.Exit(1)
		}
	}
}

// processFile обрабатывает один файл
func processFile(fname string, number, numberNonBlank, squeeze bool) error {
	var file *os.File
	var err error
	
	if fname == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(fname)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	scanner := bufio.NewScanner(file)
	lineNum := 1
	prevEmpty := false

	for scanner.Scan() {
		line := scanner.Text()
		isEmpty := len(line) == 0

		// Обработка сжатия пустых строк
		if squeeze && isEmpty && prevEmpty {
			continue
		}
		prevEmpty = isEmpty

		// Вывод номера строки
		if number {
			fmt.Printf("%6d  %s\n", lineNum, line)
			lineNum++
		} else if numberNonBlank && !isEmpty {
			fmt.Printf("%6d  %s\n", lineNum, line)
			lineNum++
		} else {
			fmt.Println(line)
		}
	}

	return scanner.Err()
}