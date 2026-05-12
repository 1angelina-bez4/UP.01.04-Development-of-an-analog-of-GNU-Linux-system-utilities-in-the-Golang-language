package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов (3 аргумента)
	var (
		help     = flag.Bool("h", false, "показать справку")
		physical = flag.Bool("P", false, "избегать символических ссылок")
		logical  = flag.Bool("L", false, "использовать логический путь (по умолчанию)")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("pwd - выводит текущую рабочую директорию")
		fmt.Println("Использование: pwd [-P] [-L]")
		fmt.Println("  -P    показывать физический путь (без символических ссылок)")
		fmt.Println("  -L    показывать логический путь (по умолчанию)")
		return
	}

	// Получение текущей директории
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: ошибка: %v\n", err)
		os.Exit(1)
	}

	// Если указан флаг -P, пытаемся получить физический путь
	if *physical {
		// Для физического пути нужно разрешить символические ссылки
		realPath, err := filepath.EvalSymlinks(dir)
		if err == nil {
			dir = realPath
		}
	}

	fmt.Println(dir)
}

// Добавляем импорт filepath
import "path/filepath"