package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileInfo структура для хранения информации о файле
type FileInfo struct {
	Name    string      // имя файла
	Size    int64       // размер в байтах
	ModTime time.Time   // время последнего изменения
	IsDir   bool        // является ли директорией
	Mode    fs.FileMode // права доступа
}

func main() {
	// Определение флагов командной строки
	var (
		showAll   = flag.Bool("a", false, "показывать скрытые файлы (начинающиеся с .)")
		longFormat = flag.Bool("l", false, "использовать длинный формат вывода")
		recursive = flag.Bool("R", false, "рекурсивно обходить поддиректории")
		reverse   = flag.Bool("r", false, "сортировать в обратном порядке")
		human     = flag.Bool("h", false, "выводить размеры в человеко-читаемом формате")
		help      = flag.Bool("help", false, "показать справку")
	)
	flag.Parse()

	// Обработка запроса справки
	if *help {
		showHelp()
		return
	}

	// Определение пути для просмотра
	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	// Запуск рекурсивного или обычного вывода
	if *recursive {
		listFilesRecursive(path, *showAll, *longFormat, *reverse, *human, "")
	} else {
		listFiles(path, *showAll, *longFormat, *reverse, *human)
	}
}

// showHelp выводит справку по использованию команды
func showHelp() {
	fmt.Println("Использование: ls [опции] [путь]")
	fmt.Println("Опции:")
	fmt.Println("  -a    показывать скрытые файлы")
	fmt.Println("  -l    подробный вывод (права, размер, время)")
	fmt.Println("  -R    рекурсивный обход директорий")
	fmt.Println("  -r    обратный порядок сортировки")
	fmt.Println("  -h    человеко-читаемые размеры (KB, MB, GB)")
	fmt.Println("  -help показать эту справку")
}

// listFiles выводит содержимое директории
func listFiles(path string, showAll, longFormat, reverse, human bool) {
	// Чтение директории
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: ошибка чтения '%s': %v\n", path, err)
		return
	}

	// Сбор информации о файлах
	var files []FileInfo
	for _, entry := range entries {
		// Пропускаем скрытые файлы, если не указан флаг -a
		if !showAll && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		
		info, err := entry.Info()
		if err != nil {
			continue // пропускаем проблемные файлы
		}
		
		files = append(files, FileInfo{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   entry.IsDir(),
			Mode:    info.Mode(),
		})
	}

	// Сортировка файлов
	sort.Slice(files, func(i, j int) bool {
		if reverse {
			return files[i].Name > files[j].Name
		}
		return files[i].Name < files[j].Name
	})

	// Вывод файлов
	for _, file := range files {
		if longFormat {
			printLongFormat(file, human)
		} else {
			printShortFormat(file, human)
		}
	}
	
	if !longFormat {
		fmt.Println() // перевод строки после короткого вывода
	}
}

// printLongFormat выводит файл в длинном формате
func printLongFormat(file FileInfo, human bool) {
	// Определение типа файла для первого символа
	fileType := "-"
	if file.IsDir {
		fileType = "d"
	}
	
	// Форматирование времени
	timeStr := file.ModTime.Format("Jan 02 15:04")
	
	// Вывод размера с возможным человеко-читаемым форматом
	sizeStr := fmt.Sprintf("%8d", file.Size)
	if human && file.Size > 0 {
		sizeStr = fmt.Sprintf("%8s", humanSize(file.Size))
	}
	
	fmt.Printf("%s%s %s %s %s\n",
		fileType,
		file.Mode.String()[1:], // убираем первый символ (тип файла)
		sizeStr,
		timeStr,
		file.Name)
}

// printShortFormat выводит файл в коротком формате
func printShortFormat(file FileInfo, human bool) {
	if human && file.Size > 0 {
		fmt.Printf("%s (%s)  ", file.Name, humanSize(file.Size))
	} else {
		fmt.Printf("%s  ", file.Name)
	}
}

// listFilesRecursive рекурсивно выводит содержимое директорий
func listFilesRecursive(path string, showAll, longFormat, reverse, human bool, prefix string) {
	fmt.Printf("%s:\n", path)
	listFiles(path, showAll, longFormat, reverse, human)
	
	// Чтение директории для рекурсивного обхода
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			// Пропускаем скрытые директории, если не указан -a
			if !showAll && strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			subPath := filepath.Join(path, entry.Name())
			fmt.Println()
			listFilesRecursive(subPath, showAll, longFormat, reverse, human, prefix+"  ")
		}
	}
}

// humanSize преобразует байты в человеко-читаемый формат
func humanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	units := []string{"K", "M", "G", "T", "P", "E"}
	return fmt.Sprintf("%.1f %sB", float64(bytes)/float64(div), units[exp])
}