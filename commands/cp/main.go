package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type CpOptions struct {
	Recursive   bool
	Force       bool
	Verbose     bool
	Preserve    bool
	Interactive bool
}

func main() {
	// Определение флагов
	var (
		help        = flag.Bool("h", false, "показать справку")
		recursive   = flag.Bool("r", false, "рекурсивное копирование")
		force       = flag.Bool("f", false, "принудительное копирование")
		verbose     = flag.Bool("v", false, "подробный вывод")
		preserve    = flag.Bool("p", false, "сохранять атрибуты")
		interactive = flag.Bool("i", false, "запрашивать подтверждение")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("cp - копирование файлов и директорий")
		fmt.Println("Использование: cp [-r] [-f] [-v] [-p] [-i] источник назначение")
		fmt.Println("  -r    рекурсивное копирование")
		fmt.Println("  -f    принудительное копирование")
		fmt.Println("  -v    подробный вывод")
		fmt.Println("  -p    сохранять атрибуты")
		fmt.Println("  -i    запрашивать подтверждение")
		return
	}

	// Проверка аргументов
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "cp: не указаны источник и назначение")
		os.Exit(1)
	}

	opts := CpOptions{
		Recursive:   *recursive,
		Force:       *force,
		Verbose:     *verbose,
		Preserve:    *preserve,
		Interactive: *interactive,
	}

	sources := flag.Args()[:len(flag.Args())-1]
	dest := flag.Args()[len(flag.Args())-1]

	// Определяем, является ли назначение директорией
	destInfo, destIsDir := os.Stat(dest)
	destIsDir = destIsDir == nil && destInfo.IsDir()

	// Копируем каждый источник
	for _, src := range sources {
		srcInfo, err := os.Stat(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cp: %v\n", err)
			continue
		}

		var destPath string
		if destIsDir {
			destPath = filepath.Join(dest, filepath.Base(src))
		} else {
			destPath = dest
		}

		if srcInfo.IsDir() {
			if !opts.Recursive {
				fmt.Fprintf(os.Stderr, "cp: пропуск директории '%s' (используйте -r)\n", src)
				continue
			}
			if err := copyDirectory(src, destPath, opts); err != nil {
				fmt.Fprintf(os.Stderr, "cp: ошибка копирования '%s': %v\n", src, err)
			}
		} else {
			if err := copyFile(src, destPath, opts); err != nil {
				fmt.Fprintf(os.Stderr, "cp: ошибка копирования '%s': %v\n", src, err)
			}
		}
	}
}

// copyFile копирует один файл
func copyFile(src, dest string, opts CpOptions) error {
	// Проверяем существование файла назначения
	if _, err := os.Stat(dest); err == nil {
		if opts.Interactive {
			var response string
			fmt.Printf("cp: перезаписать '%s'? [y/N] ", dest)
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				return nil
			}
		}
		if !opts.Force {
			return fmt.Errorf("файл '%s' уже существует", dest)
		}
	}

	// Открываем исходный файл
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Создаём файл назначения
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Копируем содержимое
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Сохраняем атрибуты если нужно
	if opts.Preserve {
		srcInfo, _ := os.Stat(src)
		os.Chmod(dest, srcInfo.Mode())
		os.Chtimes(dest, srcInfo.ModTime(), srcInfo.ModTime())
	}

	if opts.Verbose {
		fmt.Printf("cp: '%s' -> '%s'\n", src, dest)
	}

	return nil
}

// copyDirectory рекурсивно копирует директорию
func copyDirectory(src, dest string, opts CpOptions) error {
	// Создаём директорию назначения
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return err
	}

	// Читаем содержимое директории
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Копируем каждый элемент
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			if err := copyDirectory(srcPath, destPath, opts); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, destPath, opts); err != nil {
				return err
			}
		}
	}

	if opts.Verbose {
		fmt.Printf("cp: скопирована директория '%s' -> '%s'\n", src, dest)
	}

	return nil
}