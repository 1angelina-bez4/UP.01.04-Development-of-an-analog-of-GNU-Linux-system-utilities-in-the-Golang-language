package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// UnzipOptions хранит опции распаковки
type UnzipOptions struct {
	OutputDir   string
	Overwrite   bool
	Quiet       bool
	ListOnly    bool
}

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		outputDir = flag.String("d", ".", "директория для извлечения")
		overwrite = flag.Bool("o", false, "перезаписывать существующие файлы")
		quiet     = flag.Bool("q", false, "тихий режим")
		list      = flag.Bool("l", false, "только показать содержимое архива")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("unzip - распаковка ZIP архива")
		fmt.Println("Использование: unzip [-d директория] [-o] [-l] архив.zip")
		fmt.Println("  -d      директория для извлечения")
		fmt.Println("  -o      перезаписывать файлы")
		fmt.Println("  -q      тихий режим")
		fmt.Println("  -l      только показать содержимое")
		return
	}

	// Проверка наличия архива
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "unzip: нужно указать ZIP архив")
		os.Exit(1)
	}

	opts := UnzipOptions{
		OutputDir: *outputDir,
		Overwrite: *overwrite,
		Quiet:     *quiet,
		ListOnly:  *list,
	}

	zipName := flag.Arg(0)

	if opts.ListOnly {
		listZipContents(zipName, opts)
	} else {
		unzipArchive(zipName, opts)
	}
}

// unzipArchive распаковывает ZIP архив
func unzipArchive(zipName string, opts UnzipOptions) {
	// Открываем ZIP архив
	zipReader, err := zip.OpenReader(zipName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unzip: %v\n", err)
		os.Exit(1)
	}
	defer zipReader.Close()

	// Создаём выходную директорию
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "unzip: %v\n", err)
		os.Exit(1)
	}

	extracted := 0
	// Извлекаем файлы
	for _, file := range zipReader.File {
		if err := extractFile(file, opts); err != nil {
			fmt.Fprintf(os.Stderr, "unzip: ошибка извлечения %s: %v\n", file.Name, err)
		} else {
			extracted++
			if !opts.Quiet {
				fmt.Printf("unzip: извлечён '%s'\n", file.Name)
			}
		}
	}

	if !opts.Quiet {
		fmt.Printf("unzip: извлечено %d файлов\n", extracted)
	}
}

// extractFile извлекает один файл из архива
func extractFile(file *zip.File, opts UnzipOptions) error {
	destPath := filepath.Join(opts.OutputDir, file.Name)

	// Проверка на директорию
	if file.FileInfo().IsDir() {
		return os.MkdirAll(destPath, file.Mode())
	}

	// Создаём директорию для файла
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Проверка на существование
	if !opts.Overwrite {
		if _, err := os.Stat(destPath); err == nil {
			if !opts.Quiet {
				fmt.Printf("unzip: '%s' уже существует (пропуск)\n", destPath)
			}
			return nil
		}
	}

	// Открываем файл в архиве
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Создаём файл назначения
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Копируем содержимое
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Устанавливаем права доступа
	return os.Chmod(destPath, file.Mode())
}

// listZipContents показывает содержимое архива
func listZipContents(zipName string, opts UnzipOptions) {
	zipReader, err := zip.OpenReader(zipName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unzip: %v\n", err)
		os.Exit(1)
	}
	defer zipReader.Close()

	fmt.Printf("Archive: %s\n", zipName)
	fmt.Println("  Length      Date    Time    Name")
	fmt.Println("---------  ---------- -----   ----")
	
	var totalSize int64
	for _, file := range zipReader.File {
		fmt.Printf("%9d  %s  %s\n", 
			file.UncompressedSize64,
			file.Modified.Format("2006-01-02 15:04"),
			file.Name)
		totalSize += int64(file.UncompressedSize64)
	}
	
	fmt.Println("---------                     -------")
	fmt.Printf("%9d                     %d files\n", totalSize, len(zipReader.File))
}