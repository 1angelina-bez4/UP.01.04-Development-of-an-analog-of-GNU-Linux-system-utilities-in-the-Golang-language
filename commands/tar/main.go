package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TarOptions хранит опции работы с tar
type TarOptions struct {
	Create   bool
	Extract  bool
	File     string
	Verbose  bool
	Preserve bool
}

func main() {
	// Определение флагов
	var (
		help     = flag.Bool("h", false, "показать справку")
		create   = flag.Bool("c", false, "создать архив")
		extract  = flag.Bool("x", false, "извлечь архив")
		file     = flag.String("f", "", "имя архивного файла")
		verbose  = flag.Bool("v", false, "подробный вывод")
		preserve = flag.Bool("p", false, "сохранять права доступа")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("tar - создание и извлечение tar архивов")
		fmt.Println("Использование:")
		fmt.Println("  tar -c -f архив.tar файлы...  (создание)")
		fmt.Println("  tar -x -f архив.tar           (извлечение)")
		fmt.Println("Опции:")
		fmt.Println("  -c   создать архив")
		fmt.Println("  -x   извлечь архив")
		fmt.Println("  -f   файл архива")
		fmt.Println("  -v   подробный вывод")
		fmt.Println("  -p   сохранять права доступа")
		return
	}

	// Проверка обязательных опций
	if *file == "" {
		fmt.Fprintln(os.Stderr, "tar: требуется указать -f")
		os.Exit(1)
	}

	if !*create && !*extract {
		fmt.Fprintln(os.Stderr, "tar: нужно указать -c или -x")
		os.Exit(1)
	}

	opts := TarOptions{
		Create:   *create,
		Extract:  *extract,
		File:     *file,
		Verbose:  *verbose,
		Preserve: *preserve,
	}

	if opts.Create {
		if flag.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "tar: нужно указать файлы для архивации")
			os.Exit(1)
		}
		createTar(flag.Args(), opts)
	} else if opts.Extract {
		extractTar(opts)
	}
}

// createTar создаёт tar архив
func createTar(files []string, opts TarOptions) {
	// Создаём файл архива
	tarFile, err := os.Create(opts.File)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tar: %v\n", err)
		os.Exit(1)
	}
	defer tarFile.Close()

	// Создаём tar writer
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// Добавляем файлы
	for _, path := range files {
		if err := addToTar(tarWriter, path, "", opts); err != nil {
			fmt.Fprintf(os.Stderr, "tar: ошибка добавления %s: %v\n", path, err)
		}
	}
}

// addToTar добавляет файл или директорию в tar
func addToTar(tarWriter *tar.Writer, path, basePath string, opts TarOptions) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Создаём заголовок
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	
	header.Name = filepath.Join(basePath, info.Name())
	
	if opts.Verbose {
		fmt.Printf("Добавление: %s\n", header.Name)
	}

	// Записываем заголовок
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// Если это файл, копируем содержимое
	if !info.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		
		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return err
		}
	} else {
		// Рекурсивно обрабатываем директорию
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		
		for _, entry := range entries {
			newPath := filepath.Join(path, entry.Name())
			newBase := filepath.Join(basePath, info.Name(), entry.Name())
			if err := addToTar(tarWriter, newPath, newBase, opts); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// extractTar извлекает tar архив
func extractTar(opts TarOptions) {
	// Открываем архив
	tarFile, err := os.Open(opts.File)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tar: %v\n", err)
		os.Exit(1)
	}
	defer tarFile.Close()

	// Создаём tar reader
	tarReader := tar.NewReader(tarFile)

	// Извлекаем файлы
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "tar: %v\n", err)
			break
		}

		if opts.Verbose {
			fmt.Printf("Извлечение: %s\n", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(header.Name, header.FileInfo().Mode()); err != nil {
				fmt.Fprintf(os.Stderr, "tar: %v\n", err)
			}
		case tar.TypeReg:
			// Создаём директорию для файла
			if err := os.MkdirAll(filepath.Dir(header.Name), 0755); err != nil {
				fmt.Fprintf(os.Stderr, "tar: %v\n", err)
				continue
			}
			
			// Создаём файл
			file, err := os.Create(header.Name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "tar: %v\n", err)
				continue
			}
			
			// Копируем содержимое
			_, err = io.Copy(file, tarReader)
			file.Close()
			
			if err != nil {
				fmt.Fprintf(os.Stderr, "tar: %v\n", err)
				continue
			}
			
			// Устанавливаем права доступа
			if opts.Preserve {
				os.Chmod(header.Name, header.FileInfo().Mode())
			}
		}
	}
}