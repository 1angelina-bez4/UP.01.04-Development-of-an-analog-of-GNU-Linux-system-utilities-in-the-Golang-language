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
		help    = flag.Bool("h", false, "показать справку")
		lines   = flag.Int("n", 10, "количество строк")
		bytes   = flag.Int("c", 0, "количество байт")
		quiet   = flag.Bool("q", false, "не выводить заголовки файлов")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("head - вывод первых строк файла")
		fmt.Println("Использование: head [-n строки] [-c байты] [-q] [файлы...]")
		fmt.Println("  -n   количество строк (по умолчанию 10)")
		fmt.Println("  -c   количество байт")
		fmt.Println("  -q   не выводить заголовки")
		return
	}

	// Проверка аргументов
	if *bytes < 0 {
		fmt.Fprintln(os.Stderr, "head: количество байт не может быть отрицательным")
		os.Exit(1)
	}
	if *lines < 0 {
		fmt.Fprintln(os.Stderr, "head: количество строк не может быть отрицательным")
		os.Exit(1)
	}
	if *bytes > 0 && *lines > 0 {
		fmt.Fprintln(os.Stderr, "head: нельзя использовать -n и -c одновременно")
		os.Exit(1)
	}

	// Если файлы не указаны, читаем stdin
	if flag.NArg() == 0 {
		if *bytes > 0 {
			headBytes(os.Stdin, *bytes)
		} else {
			headLines(os.Stdin, *lines)
		}
		return
	}

	// Обрабатываем каждый файл
	for i, fname := range flag.Args() {
		file, err := os.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "head: %v\n", err)
			continue
		}
		
		// Выводим заголовок
		if len(flag.Args()) > 1 && !*quiet {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("==> %s <==\n", fname)
		}
		
		if *bytes > 0 {
			headBytes(file, *bytes)
		} else {
			headLines(file, *lines)
		}
		file.Close()
	}
}

// headLines выводит первые N строк
func headLines(file *os.File, n int) {
	scanner := bufio.NewScanner(file)
	for i := 0; i < n && scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "head: ошибка чтения: %v\n", err)
	}
}

// headBytes выводит первые N байт
func headBytes(file *os.File, n int) {
	buf := make([]byte, n)
	count, err := file.Read(buf)
	if err != nil && err.Error() != "EOF" {
		fmt.Fprintf(os.Stderr, "head: ошибка чтения: %v\n", err)
		return
	}
	os.Stdout.Write(buf[:count])
}