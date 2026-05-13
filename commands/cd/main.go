package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// Флаги
	recursive := flag.Bool("r", false, "рекурсивное копирование")
	force := flag.Bool("f", false, "принудительное копирование")
	verbose := flag.Bool("v", false, "подробный вывод")
	help := flag.Bool("h", false, "справка")
	flag.Parse()

	if *help {
		fmt.Println("cp [-r] [-f] [-v] источник назначение")
		return
	}

	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "cp: нужен источник и назначение")
		os.Exit(1)
	}

	args := flag.Args()
	sources := args[:len(args)-1]
	dest := args[len(args)-1]

	// Проверяем, директория ли назначение
	destInfo, _ := os.Stat(dest)
	isDir := destInfo != nil && destInfo.IsDir()

	// Копируем
	for _, src := range sources {
		srcInfo, err := os.Stat(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cp: %v\n", err)
			continue
		}

		// Определяем путь назначения
		destPath := dest
		if isDir {
			destPath = filepath.Join(dest, filepath.Base(src))
		}

		// Копируем файл или директорию
		if srcInfo.IsDir() {
			if !*recursive {
				fmt.Fprintf(os.Stderr, "cp: '%s' - директория (нужен -r)\n", src)
				continue
			}
			copyDir(src, destPath, *force, *verbose)
		} else {
			copyFile(src, destPath, *force, *verbose)
		}
	}
}

func copyFile(src, dst string, force, verbose bool) error {
	// Проверяем существование
	if _, err := os.Stat(dst); err == nil && !force {
		return fmt.Errorf("'%s' уже существует", dst)
	}

	// Открываем исходный
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Создаем назначение
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Копируем
	_, err = io.Copy(out, in)
	if err == nil && verbose {
		fmt.Printf("%s -> %s\n", src, dst)
	}
	return err
}

func copyDir(src, dst string, force, verbose bool) error {
	// Создаем директорию
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Читаем содержимое
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Копируем каждый элемент
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			copyDir(srcPath, dstPath, force, verbose)
		} else {
			copyFile(srcPath, dstPath, force, verbose)
		}
	}

	if verbose {
		fmt.Printf("dir: %s -> %s\n", src, dst)
	}
	return nil
}
