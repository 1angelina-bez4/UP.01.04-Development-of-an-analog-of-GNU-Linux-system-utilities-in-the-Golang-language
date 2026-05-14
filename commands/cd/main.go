package main

import (
	"flag"
	"fmt"
	"os"
)

// Глобальная переменная для хранения предыдущей директории
var prevDir string

func main() {
	//Создание переменных флагов с выводом информации
	var (
		help    = flag.Bool("h", false, "показать справку")
		toPrev  = flag.Bool("prev", false, "перейти в предыдущую директорию (cd -)")
		verbose = flag.Bool("v", false, "выводить путь после перехода")
	)
	flag.Parse()

	// Вывод справки с помощью флага -h
	if *help {
		fmt.Println("cd - смена текущей директории")
		fmt.Println("Использование: cd [директория]")
		fmt.Println("  cd -       перейти в предыдущую директорию")
		fmt.Println("  cd -prev   перейти в предыдущую директорию")
		fmt.Println("  -v         выводить путь после перехода")
		return
	}

	// Получаем текущую директорию
	current, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %v\n", err)
		os.Exit(1)
	}

	// Обработка перехода в предыдущую директорию
	if *toPrev || (flag.NArg() > 0 && flag.Arg(0) == "-") {
		if prevDir == "" {
			fmt.Fprintln(os.Stderr, "cd: нет предыдущей директории")
			os.Exit(1)
		}
		
		// Меняем директорию
		if err := os.Chdir(prevDir); err != nil {
			fmt.Fprintf(os.Stderr, "cd: %v\n", err)
			os.Exit(1)
		}
		
		// Обновляем prevDir
		prevDir = current
		
		if *verbose {
			newDir, _ := os.Getwd()
			fmt.Println(newDir)
		}
		return
	}

	// Определение целевой директории
	target := "."
	if flag.NArg() > 0 {
		target = flag.Arg(0)
	}

	// Проверка существования директории
	info, err := os.Stat(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %v\n", err)
		os.Exit(1)
	}
	
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "cd: %s: не является директорией\n", target)
		os.Exit(1)
	}

	// Сохраняем текущую директорию перед переходом
	prevDir = current

	// Выполняем переход
	if err := os.Chdir(target); err != nil {
		fmt.Fprintf(os.Stderr, "cd: %v\n", err)
		os.Exit(1)
	}

	if *verbose {
		newDir, _ := os.Getwd()
		fmt.Println(newDir)
	}
}
