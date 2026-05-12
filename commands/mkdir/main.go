package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		parents   = flag.Bool("p", false, "создавать промежуточные директории")
		verbose   = flag.Bool("v", false, "выводить сообщения о создании")
		mode      = flag.String("m", "0755", "права доступа (восьмеричные)")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("mkdir - создание директорий")
		fmt.Println("Использование: mkdir [-p] [-m режим] директория...")
		fmt.Println("  -p    создавать промежуточные директории")
		fmt.Println("  -m    установить права доступа")
		fmt.Println("  -v    выводить подробную информацию")
		return
	}

	// Проверка наличия аргументов
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "mkdir: не указана директория")
		os.Exit(1)
	}

	// Парсим права доступа
	var perm os.FileMode
	_, err := fmt.Sscanf(*mode, "%o", &perm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: неверный режим '%s'\n", *mode)
		os.Exit(1)
	}

	// Создаём директории
	for _, dir := range flag.Args() {
		var err error
		if *parents {
			err = os.MkdirAll(dir, perm)
		} else {
			err = os.Mkdir(dir, perm)
		}
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
			os.Exit(1)
		}
		
		if *verbose {
			fmt.Printf("mkdir: создана директория '%s'\n", dir)
		}
	}
}