package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	d := flag.String("d", ".", "директория")
	o := flag.Bool("o", false, "перезапись")
	q := flag.Bool("q", false, "тихо")
	l := flag.Bool("l", false, "список")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "unzip: нужен архив")
		os.Exit(1)
	}

	zipName := flag.Arg(0)

	if *l {
		list(zipName)
		return
	}

	// Открываем архив
	r, err := zip.OpenReader(zipName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unzip: %v\n", err)
		os.Exit(1)
	}
	defer r.Close()

	// Создаем директорию
	os.MkdirAll(*d, 0755)

	count := 0
	for _, f := range r.File {
		path := filepath.Join(*d, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		// Создаем директорию для файла
		os.MkdirAll(filepath.Dir(path), 0755)

		// Проверяем существование
		if !*o {
			if _, err := os.Stat(path); err == nil {
				if !*q {
					fmt.Printf("unzip: пропущен %s\n", f.Name)
				}
				continue
			}
		}

		// Распаковываем
		src, _ := f.Open()
		dst, _ := os.Create(path)
		io.Copy(dst, src)
		src.Close()
		dst.Close()
		os.Chmod(path, f.Mode())
		count++

		if !*q {
			fmt.Printf("unzip: %s\n", f.Name)
		}
	}

	if !*q {
		fmt.Printf("unzip: %d файлов\n", count)
	}
}

func list(zipName string) {
	r, err := zip.OpenReader(zipName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unzip: %v\n", err)
		return
	}
	defer r.Close()

	fmt.Printf("Архив: %s\n", zipName)
	fmt.Println("  Размер  Имя")
	fmt.Println("  ------  ----")

	var total int64
	for _, f := range r.File {
		fmt.Printf("%8d  %s\n", f.UncompressedSize64, f.Name)
		total += int64(f.UncompressedSize64)
	}
	fmt.Printf("  ------  ----\n")
	fmt.Printf("%8d  (%d файлов)\n", total, len(r.File))
}
