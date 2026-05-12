package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FindOptions struct {
	Name      string
	Type      string
	Size      string
	MaxDepth  int
	Exec      string
}

func main() {
	// Определение флагов
	var (
		help     = flag.Bool("h", false, "показать справку")
		name     = flag.String("name", "", "имя файла (возможен шаблон)")
		fileType = flag.String("type", "", "тип: f (файл), d (директория)")
		size     = flag.String("size", "", "размер: +10M, -1G, 500k")
		maxDepth = flag.Int("maxdepth", -1, "максимальная глубина")
		exec     = flag.String("exec", "", "выполнить команду для каждого найденного файла")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("find - поиск файлов")
		fmt.Println("Использование: find [путь] [-name шаблон] [-type тип] [-size размер] [-maxdepth глубина] [-exec команда]")
		fmt.Println("  -name      имя файла (поддерживает * и ?)")
		fmt.Println("  -type      тип: f (файл), d (директория)")
		fmt.Println("  -size      размер: +10M, -1G, 500k (b, k, M, G)")
		fmt.Println("  -maxdepth  максимальная глубина поиска")
		fmt.Println("  -exec      выполнить команду ({} заменяется на путь к файлу)")
		return
	}

	// Определяем путь для поиска
	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	opts := FindOptions{
		Name:     *name,
		Type:     *fileType,
		Size:     *size,
		MaxDepth: *maxDepth,
		Exec:     *exec,
	}

	// Выполняем поиск
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // пропускаем ошибки доступа
		}

		// Проверяем глубину
		if opts.MaxDepth >= 0 {
			depth := strings.Count(currentPath[len(path):], string(os.PathSeparator))
			if depth > opts.MaxDepth {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Пропускаем сам начальный путь
		if currentPath == path {
			return nil
		}

		// Проверка типа файла
		if opts.Type != "" {
			switch opts.Type {
			case "f":
				if info.IsDir() {
					return nil
				}
			case "d":
				if !info.IsDir() {
					return nil
				}
			default:
				return nil
			}
		}

		// Проверка имени файла
		if opts.Name != "" {
			matched, err := filepath.Match(opts.Name, info.Name())
			if err != nil || !matched {
				return nil
			}
		}

		// Проверка размера
		if opts.Size != "" {
			if !checkSize(info.Size(), opts.Size) {
				return nil
			}
		}

		// Вывод результата
		if opts.Exec != "" {
			executeCommand(opts.Exec, currentPath)
		} else {
			fmt.Println(currentPath)
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "find: %v\n", err)
		os.Exit(1)
	}
}

// checkSize проверяет, соответствует ли размер файла условию
func checkSize(fileSize int64, sizePattern string) bool {
	var multiplier int64 = 1
	var compareOp string
	
	// Определяем оператор сравнения
	if strings.HasPrefix(sizePattern, "+") {
		compareOp = ">"
		sizePattern = sizePattern[1:]
	} else if strings.HasPrefix(sizePattern, "-") {
		compareOp = "<"
		sizePattern = sizePattern[1:]
	} else {
		compareOp = "="
	}
	
	// Определяем множитель
	if strings.HasSuffix(sizePattern, "b") {
		multiplier = 1
		sizePattern = sizePattern[:len(sizePattern)-1]
	} else if strings.HasSuffix(sizePattern, "k") {
		multiplier = 1024
		sizePattern = sizePattern[:len(sizePattern)-1]
	} else if strings.HasSuffix(sizePattern, "M") {
		multiplier = 1024 * 1024
		sizePattern = sizePattern[:len(sizePattern)-1]
	} else if strings.HasSuffix(sizePattern, "G") {
		multiplier = 1024 * 1024 * 1024
		sizePattern = sizePattern[:len(sizePattern)-1]
	}
	
	var targetSize int64
	fmt.Sscanf(sizePattern, "%d", &targetSize)
	targetSize *= multiplier
	
	switch compareOp {
	case ">":
		return fileSize > targetSize
	case "<":
		return fileSize < targetSize
	default:
		return fileSize == targetSize
	}
}

// executeCommand выполняет команду для найденного файла
func executeCommand(command, filePath string) {
	// Заменяем {} на путь к файлу
	cmdStr := strings.ReplaceAll(command, "{}", filePath)
	
	// Простой вывод команды (в реальной оболочке нужно выполнять)
	fmt.Printf("Выполнение: %s\n", cmdStr)
	
	// Здесь можно добавить реальное выполнение команды
	// parts := strings.Fields(cmdStr)
	// cmd := exec.Command(parts[0], parts[1:]...)
	// cmd.Run()
}