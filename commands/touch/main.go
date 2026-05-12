package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		noCreate  = flag.Bool("c", false, "не создавать файл, если не существует")
		date      = flag.String("t", "", "установить время (формат: YYYYMMDDHHMM.SS)")
		reference = flag.String("r", "", "использовать время указанного файла")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("touch - изменение временных меток файла")
		fmt.Println("Использование: touch [-c] [-t время] [-r файл] файл...")
		fmt.Println("  -c      не создавать новый файл")
		fmt.Println("  -t      установить указанное время")
		fmt.Println("  -r      использовать время указанного файла")
		return
	}

	// Проверка наличия файлов
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "touch: не указан файл")
		os.Exit(1)
	}

	// Определяем время для установки
	var modTime time.Time
	var err error
	
	if *reference != "" {
		modTime, err = getReferenceTime(*reference)
		if err != nil {
			fmt.Fprintf(os.Stderr, "touch: %v\n", err)
			os.Exit(1)
		}
	} else if *date != "" {
		modTime, err = parseTime(*date)
		if err != nil {
			fmt.Fprintf(os.Stderr, "touch: %v\n", err)
			os.Exit(1)
		}
	} else {
		modTime = time.Now()
	}

	// Обрабатываем каждый файл
	for _, fname := range flag.Args() {
		if err := touchFile(fname, modTime, *noCreate); err != nil {
			fmt.Fprintf(os.Stderr, "touch: %v\n", err)
			os.Exit(1)
		}
	}
}

// touchFile обрабатывает один файл
func touchFile(fname string, modTime time.Time, noCreate bool) error {
	// Проверяем существование файла
	_, err := os.Stat(fname)
	if os.IsNotExist(err) {
		if noCreate {
			return nil
		}
		// Создаём файл
		file, err := os.Create(fname)
		if err != nil {
			return err
		}
		file.Close()
	} else if err != nil {
		return err
	}
	
	// Изменяем временные метки
	return os.Chtimes(fname, modTime, modTime)
}

// parseTime парсит время из строки
func parseTime(timeStr string) (time.Time, error) {
	// Формат: 202312151430.45
	if len(timeStr) != 15 {
		return time.Time{}, fmt.Errorf("неверный формат, используйте YYYYMMDDHHMM.SS")
	}
	
	year := timeStr[0:4]
	month := timeStr[4:6]
	day := timeStr[6:8]
	hour := timeStr[8:10]
	min := timeStr[10:12]
	sec := timeStr[13:15]
	
	layout := "20060102150405"
	timeStrWithoutDot := year + month + day + hour + min + sec
	
	return time.ParseInLocation(layout, timeStrWithoutDot, time.Local)
}

// getReferenceTime получает время из файла
func getReferenceTime(fname string) (time.Time, error) {
	info, err := os.Stat(fname)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}