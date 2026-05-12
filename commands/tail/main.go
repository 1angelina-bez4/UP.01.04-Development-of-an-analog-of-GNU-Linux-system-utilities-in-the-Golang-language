package main

import (
	"bufio"
	"container/list"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов
	var (
		help    = flag.Bool("h", false, "показать справку")
		lines   = flag.Int("n", 10, "количество строк")
		follow  = flag.Bool("f", false, "следить за файлом")
		quiet   = flag.Bool("q", false, "не выводить заголовки")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("tail - вывод последних строк файла")
		fmt.Println("Использование: tail [-n строки] [-f] [-q] [файлы...]")
		fmt.Println("  -n    количество строк (по умолчанию 10)")
		fmt.Println("  -f    следить за файлом")
		fmt.Println("  -q    не выводить заголовки")
		return
	}

	if *lines < 0 {
		fmt.Fprintln(os.Stderr, "tail: количество строк не может быть отрицательным")
		os.Exit(1)
	}

	// Если файлы не указаны, читаем stdin
	if flag.NArg() == 0 {
		tailLines(os.Stdin, *lines)
		return
	}

	// Обрабатываем каждый файл
	for i, fname := range flag.Args() {
		file, err := os.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "tail: %v\n", err)
			continue
		}

		// Выводим заголовок
		if len(flag.Args()) > 1 && !*quiet {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("==> %s <==\n", fname)
		}

		tailLines(file, *lines)
		file.Close()

		// Если нужно следить за файлом
		if *follow {
			followFile(fname)
		}
	}
}

// tailLines выводит последние N строк файла
func tailLines(file *os.File, n int) {
	scanner := bufio.NewScanner(file)
	
	// Используем список для хранения последних N строк
	linesList := list.New()
	
	for scanner.Scan() {
		line := scanner.Text()
		
		if linesList.Len() >= n {
			linesList.Remove(linesList.Front())
		}
		linesList.PushBack(line)
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "tail: ошибка чтения: %v\n", err)
		return
	}
	
	// Выводим строки
	for e := linesList.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(string))
	}
}

// followFile следит за файлом и выводит новые строки
func followFile(fname string) {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tail: %v\n", err)
		return
	}
	defer file.Close()
	
	// Перемещаемся в конец файла
	file.Seek(0, 2)
	
	scanner := bufio.NewScanner(file)
	for {
		if scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}
}