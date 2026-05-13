package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FindOptions struct {
	Name     string
	Type     string
	Size     string
	MaxDepth int
	Exec     string
}

func main() {
	var (
		help     = flag.Bool("h", false, "показать справку")
		name     = flag.String("name", "", "имя файла (возможен шаблон)")
		fileType = flag.String("type", "", "тип: f (файл), d (директория)")
		size     = flag.String("size", "", "размер: +10M, -1G, 500k")
		maxDepth = flag.Int("maxdepth", -1, "максимальная глубина")
		exec     = flag.String("exec", "", "выполнить команду для каждого найденного файла")
	)
	flag.Parse()

	if *help {
		fmt.Println("find - поиск файлов")
		fmt.Println("Использование: find [путь] [-name шаблон] [-type тип] [-size размер] [-maxdepth глубина] [-exec команда]")
		return
	}

	// Путь поиска
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	opts := FindOptions{
		Name:     *name,
		Type:     *fileType,
		Size:     *size,
		MaxDepth: *maxDepth,
		Exec:     *exec,
	}

	// Нормализуем путь
	root, _ = filepath.Abs(root)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // пропускаем ошибки
		}

		// Вычисляем глубину
		if opts.MaxDepth >= 0 {
			relPath, _ := filepath.Rel(root, path)
			depth := 0
			if relPath != "." {
				depth = strings.Count(relPath, string(os.PathSeparator)) + 1
			}
			if depth > opts.MaxDepth {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Проверяем тип
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

		// Проверяем имя
		if opts.Name != "" {
			matched, err := filepath.Match(opts.Name, info.Name())
			if err != nil || !matched {
				return nil
			}
		}

		// Проверяем размер
		if opts.Size != "" {
			if !checkSize(info.Size(), opts.Size) {
				return nil
			}
		}

		// Выводим результат
		if opts.Exec != "" {
			executeCommand(opts.Exec, path)
		} else {
			fmt.Println(path)
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "find: %v\n", err)
	}
}

func checkSize(fileSize int64, pattern string) bool {
	var multiplier int64 = 1
	var op string

	if strings.HasPrefix(pattern, "+") {
		op = ">"
		pattern = pattern[1:]
	} else if strings.HasPrefix(pattern, "-") {
		op = "<"
		pattern = pattern[1:]
	} else {
		op = "="
	}

	if strings.HasSuffix(pattern, "b") {
		multiplier = 1
		pattern = pattern[:len(pattern)-1]
	} else if strings.HasSuffix(pattern, "k") {
		multiplier = 1024
		pattern = pattern[:len(pattern)-1]
	} else if strings.HasSuffix(pattern, "M") {
		multiplier = 1024 * 1024
		pattern = pattern[:len(pattern)-1]
	} else if strings.HasSuffix(pattern, "G") {
		multiplier = 1024 * 1024 * 1024
		pattern = pattern[:len(pattern)-1]
	}

	var size int64
	fmt.Sscanf(pattern, "%d", &size)
	size *= multiplier

	switch op {
	case ">":
		return fileSize > size
	case "<":
		return fileSize < size
	default:
		return fileSize == size
	}
}

func executeCommand(command, filePath string) {
	cmdStr := strings.ReplaceAll(command, "{}", filePath)

	// Разбиваем команду на части
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return
	}

	// Выполняем команду
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "find: ошибка выполнения '%s': %v\n", cmdStr, err)
	}
}
