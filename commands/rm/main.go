package main

import (
	"flag"
	"fmt"
	"os"
)

// RmOptions хранит опции удаления
type RmOptions struct {
	Recursive bool
	Force     bool
	Verbose   bool
	Interactive bool
}

func main() {
	// Определение флагов
	var (
		help        = flag.Bool("h", false, "показать справку")
		recursive   = flag.Bool("R", false, "рекурсивное удаление")
		force       = flag.Bool("f", false, "игнорировать ошибки")
		verbose     = flag.Bool("v", false, "выводить информацию")
		interactive = flag.Bool("i", false, "запрашивать подтверждение")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("rm - удаление файлов и директорий")
		fmt.Println("Использование: rm [-R] [-f] [-v] [-i] файлы...")
		fmt.Println("  -R   рекурсивное удаление")
		fmt.Println("  -f   игнорировать ошибки")
		fmt.Println("  -v   подробный вывод")
		fmt.Println("  -i   запрашивать подтверждение")
		return
	}

	opts := RmOptions{
		Recursive:   *recursive,
		Force:       *force,
		Verbose:     *verbose,
		Interactive: *interactive,
	}

	// Проверка наличия аргументов
	if flag.NArg() == 0 && !opts.Force {
		fmt.Fprintln(os.Stderr, "rm: не указан файл")
		os.Exit(1)
	}

	// Удаляем файлы
	for _, path := range flag.Args() {
		if err := removeFile(path, opts); err != nil {
			if !opts.Force {
				fmt.Fprintf(os.Stderr, "rm: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

// removeFile удаляет файл или директорию
func removeFile(path string, opts RmOptions) error {
	// Проверяем существование
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Запрашиваем подтверждение
	if opts.Interactive {
		var response string
		fmt.Printf("rm: удалить '%s'? [y/N] ", path)
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return nil
		}
	}

	// Удаляем
	if info.IsDir() && !opts.Recursive {
		// Проверяем, пуста ли директория
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		if len(entries) > 0 {
			return fmt.Errorf("%s: директория не пуста, используйте -R", path)
		}
		return os.Remove(path)
	}

	if info.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}

	if err == nil && opts.Verbose {
		fmt.Printf("rm: удалён '%s'\n", path)
	}
	return err
}
