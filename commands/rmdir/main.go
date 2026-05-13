package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Определение флагов
	var (
		help       = flag.Bool("h", false, "показать справку")
		ignore     = flag.Bool("i", false, "игнорировать ошибки для непустых директорий")
		parents    = flag.Bool("p", false, "удалять родительские директории")
		verbose    = flag.Bool("v", false, "выводить сообщения")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("rmdir - удаление пустых директорий")
		fmt.Println("Использование: rmdir [-p] [-v] директория...")
		fmt.Println("  -p      удалять родительские директории")
		fmt.Println("  -v      выводить подробную информацию")
		fmt.Println("  -i  игнорировать ошибки")
		return
	}

	// Проверка наличия аргументов
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "rmdir: не указана директория")
		os.Exit(1)
	}

	// Удаляем директории
	for _, dir := range flag.Args() {
		if *parents {
			// Удаляем путь целиком (как rmdir -p)
			err := removePathParents(dir, *verbose, *ignore)
			if err != nil && !*ignore {
				fmt.Fprintf(os.Stderr, "rmdir: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Обычное удаление
			err := os.Remove(dir)
			if err != nil {
				if *ignore {
					continue
				}
				fmt.Fprintf(os.Stderr, "rmdir: %v\n", err)
				os.Exit(1)
			}
			if *verbose {
				fmt.Printf("rmdir: удалена '%s'\n", dir)
			}
		}
	}
}

// removePathParents удаляет директорию и её пустых родителей
func removePathParents(path string, verbose, ignore bool) error {
	for {
		// Пытаемся удалить текущую директорию
		err := os.Remove(path)
		if err != nil {
			if ignore {
				return nil
			}
			return err
		}
		
		if verbose {
			fmt.Printf("rmdir: удалена '%s'\n", path)
		}
		
		// Переходим к родительской директории
		parent := getParentDir(path)
		if parent == path || parent == "" || parent == "/" {
			break
		}
		path = parent
	}
	return nil
}

// getParentDir возвращает родительскую директорию
func getParentDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[:i]
		}
	}
	return ""
}
