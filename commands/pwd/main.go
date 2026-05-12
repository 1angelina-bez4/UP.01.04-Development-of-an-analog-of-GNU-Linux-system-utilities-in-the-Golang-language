package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	//Создаем  переменные для отображения информации при ошибках
	var (
		help     = flag.Bool("h", false, "показать справку")
		physical = flag.Bool("P", false, "избегать символических ссылок")
	)

	//Вызов флагов 
	flag.Parse()

	//Вывод информации при указании флага -h
	if *help {
		fmt.Println("pwd - выводит текущую рабочую директорию")
		fmt.Println("Использование: pwd [-P] [-L]")
		fmt.Println("  -P    показывать физический путь (без символических ссылок)")
		fmt.Println("  -L    показывать логический путь (по умолчанию)")
		return
	}

	// Если указан флаг -P, пытаемся получить физический путь
	if *physical {
		// Для физического пути нужно разрешить символические ссылки
		realPath, err := filepath.EvalSymlinks(dir)
		if err == nil {
			dir = realPath
		}
	}

	// Получение текущей директории
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pwd: ошибка: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(dir)
}
