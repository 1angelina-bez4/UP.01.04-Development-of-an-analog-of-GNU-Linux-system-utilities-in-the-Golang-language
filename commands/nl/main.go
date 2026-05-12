package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// NumberingStyle определяет стиль нумерации
type NumberingStyle int

const (
	NumberAll NumberingStyle = iota
	NumberNonEmpty
	NumberNone
)

func main() {
	// Определение флагов
	var (
		help        = flag.Bool("h", false, "показать справку")
		bodyStyle   = flag.String("b", "t", "стиль нумерации: a (все), t (непустые), n (нет)")
		numberWidth = flag.Int("w", 6, "ширина номера строки")
		separator   = flag.String("s", "\t", "разделитель между номером и текстом")
		startNum    = flag.Int("v", 1, "начальный номер")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("nl - нумерация строк")
		fmt.Println("Использование: nl [-b стиль] [-w ширина] [-s разделитель] [-v старт] [файл]")
		fmt.Println("  -b стиль   a=все, t=только непустые, n=нет")
		fmt.Println("  -w ширина  ширина поля номера")
		fmt.Println("  -s строка  разделитель")
		fmt.Println("  -v число   начальное значение")
		return
	}

	// Определяем стиль нумерации
	var style NumberingStyle
	switch *bodyStyle {
	case "a":
		style = NumberAll
	case "t":
		style = NumberNonEmpty
	case "n":
		style = NumberNone
	default:
		fmt.Fprintf(os.Stderr, "nl: неверный стиль: %s\n", *bodyStyle)
		os.Exit(1)
	}

	// Проверка ширины
	if *numberWidth < 1 || *numberWidth > 20 {
		fmt.Fprintf(os.Stderr, "nl: ширина должна быть от 1 до 20\n")
		os.Exit(1)
	}

	// Открываем файл или используем stdin
	var file *os.File
	if flag.NArg() > 0 {
		var err error
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "nl: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
	} else {
		file = os.Stdin
	}

	// Обработка строк
	scanner := bufio.NewScanner(file)
	lineNum := *startNum
	
	for scanner.Scan() {
		line := scanner.Text()
		shouldNumber := false

		switch style {
		case NumberAll:
			shouldNumber = true
		case NumberNonEmpty:
			if strings.TrimSpace(line) != "" {
				shouldNumber = true
			}
		case NumberNone:
			shouldNumber = false
		}

		if shouldNumber {
			fmt.Printf("%*d%s%s\n", *numberWidth, lineNum, *separator, line)
			lineNum++
		} else {
			fmt.Printf("%*s%s%s\n", *numberWidth, "", *separator, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "nl: ошибка чтения: %v\n", err)
		os.Exit(1)
	}
}