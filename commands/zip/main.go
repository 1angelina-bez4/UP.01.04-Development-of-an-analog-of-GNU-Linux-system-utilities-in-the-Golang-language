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

// ZipOptions хранит опции архивации
type ZipOptions struct {
	Recursive   bool
	Compression int
	Quiet       bool
	Exclude     []string
}

func main() {
	// Определение флагов
	var (
		help       = flag.Bool("h", false, "показать справку")
		recursive  = flag.Bool("r", false, "рекурсивное добавление")
		level      = flag.Int("l", 6, "уровень сжатия (0-9)")
		quiet      = flag.Bool("q", false, "не выводить информацию")
		exclude    = flag.String("x", "", "исключить файлы (через запятую)")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("zip - создание ZIP архива")
		fmt.Println("Использование: zip [-r] [-l уровень] [-x шаблон] архив.zip файлы...")
		fmt.Println("  -r      рекурсивно добавлять директории")
		fmt.Println("  -l      уровень сжатия (0-9)")
		fmt.Println("  -q      тихий режим")
		fmt.Println("  -x      исключить файлы (шаблоны через запятую)")
		return
	}

	// Проверка уровня сжатия
	if *level < 0 || *level > 9 {
		fmt.Fprintln(os.Stderr, "zip: уровень сжатия должен быть от 0 до 9")
		os.Exit(1)
	}

	// Проверка наличия аргументов
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "zip: нужно указать архив и файлы")
		os.Exit(1)
	}

	// Парсим исключения
	excludeList := []string{}
	if *exclude != "" {
		excludeList = strings.Split(*exclude, ",")
	}

	options := ZipOptions{
		Recursive:   *recursive,
		Compression: *level,
		Quiet:       *quiet,
		Exclude:     excludeList,
	}

	zipName := flag.Arg(0)
	files := flag.Args()[1:]

	// Создаём ZIP архив
	if err := createZip(zipName, files, options); err != nil {
		fmt.Fprintf(os.Stderr, "zip: %v\n", err)
		os.Exit(1)
	}

	if !options.Quiet {
		fmt.Printf("zip: создан архив '%s'\n", zipName)
	}
}

// createZip создаёт ZIP архив
func createZip(zipName string, files []string, opts ZipOptions) error {
	// Создаём файл архива
	zipFile, err := os.Create(zipName)
	if err != nil {
		return err
		defer zipFile.Close()

	// Создаём ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Добавляем файлы в архив
	for _, fname := range files {
		if err := addToZip(zipWriter, fname, "", opts); err != nil {
			return fmt.Errorf("ошибка добавления %s: %v", fname, err)
		}
	}

	return nil
}

// addToZip добавляет файл или директорию в ZIP архив
func addToZip(zipWriter *zip.Writer, path, basePath string, opts ZipOptions) error {
	// Проверяем, нужно ли исключить файл
	if shouldExclude(path, opts.Exclude) {
		if !opts.Quiet {
			fmt.Printf("zip: исключён '%s'\n", path)
		}
		return nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Обработка директории
	if info.IsDir() {
		if !opts.Recursive {
			return fmt.Errorf("'%s' - директория, используйте -r", path)
		}
		
		// Рекурсивно обходим директорию
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		
		for _, entry := range entries {
			newPath := filepath.Join(path, entry.Name())
			newBase := filepath.Join(basePath, entry.Name())
			if err := addToZip(zipWriter, newPath, newBase, opts); err != nil {
				return err
			}
		}
		return nil
	}

	// Добавляем файл
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Определяем имя в архиве
	zipPath := basePath
	if zipPath == "" {
		zipPath = filepath.Base(path)
	}

	// Создаём запись в архиве
	writer, err := zipWriter.Create(zipPath)
	if err != nil {
		return err
	}

	// Копируем содержимое
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	if !opts.Quiet {
		fmt.Printf("zip: добавлен '%s'\n", path)
	}
	return nil
}

// shouldExclude проверяет, нужно ли исключить файл
func shouldExclude(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}