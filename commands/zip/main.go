package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	r := flag.Bool("r", false, "рекурсивно")
	l := flag.Int("l", 6, "сжатие 0-9")
	q := flag.Bool("q", false, "тихо")
	x := flag.String("x", "", "исключить")
	L := flag.Bool("L", false, "список")
	flag.Parse()

	if *L {
		list(flag.Arg(0))
		return
	}

	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "zip: нужен архив и файлы")
		os.Exit(1)
	}

	zipName, files := flag.Arg(0), flag.Args()[1:]
	excludes := strings.Split(*x, ",")

	// Создаем архив
	f, _ := os.Create(zipName)
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()

	for _, file := range files {
		add(w, file, "", *r, *l, *q, excludes)
	}
	if !*q {
		fmt.Printf("zip: создан %s\n", zipName)
	}
}

func add(w *zip.Writer, path, base string, rec bool, lvl int, quiet bool, exclude []string) {
	// Проверка исключения
	for _, pat := range exclude {
		if ok, _ := filepath.Match(pat, filepath.Base(path)); ok {
			return
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.IsDir() {
		if !rec {
			fmt.Fprintf(os.Stderr, "zip: %s - директория, нужен -r\n", path)
			return
		}
		entries, _ := os.ReadDir(path)
		for _, e := range entries {
			add(w, filepath.Join(path, e.Name()), filepath.Join(base, e.Name()), rec, lvl, quiet, exclude)
		}
		return
	}

	// Добавляем файл
	file, _ := os.Open(path)
	defer file.Close()

	// Заголовок
	header := &zip.FileHeader{Name: base, Method: zip.Deflate}
	if base == "" {
		header.Name = filepath.Base(path)
	}
	if lvl == 0 {
		header.Method = zip.Store
	}

	wr, _ := w.CreateHeader(header)
	io.Copy(wr, file)

	if !quiet {
		fmt.Printf("zip: %s -> %s\n", path, header.Name)
	}
}

func list(name string) {
	r, err := zip.OpenReader(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "zip: %v\n", err)
		return
	}
	defer r.Close()

	fmt.Printf("Архив: %s\n", name)
	for _, f := range r.File {
		fmt.Printf("  %8d  %s\n", f.UncompressedSize64, f.Name)
	}
}
