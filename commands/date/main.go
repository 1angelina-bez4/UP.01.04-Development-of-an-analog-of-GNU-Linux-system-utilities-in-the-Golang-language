package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		utc       = flag.Bool("u", false, "показать UTC время")
		format    = flag.String("format", "default", "формат вывода: default, rfc3339, unix, kitchen")
		dateStr   = flag.String("d", "", "показать указанную дату")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("date - вывод даты и времени")
		fmt.Println("Использование: date [-u] [-format формат] [-d дата]")
		fmt.Println("  -u        показать UTC время")
		fmt.Println("  -format   формат: default, rfc3339, unix, kitchen")
		fmt.Println("  -d        показать указанную дату")
		return
	}

	var now time.Time
	
	// Если указана конкретная дата
	if *dateStr != "" {
		layouts := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			"02.01.2006",
			"02 Jan 2006",
		}
		
		parsed := false
		for _, layout := range layouts {
			t, err := time.Parse(layout, *dateStr)
			if err == nil {
				now = t
				parsed = true
				break
			}
		}
		
		if !parsed {
			fmt.Fprintf(nil, "date: не удалось распознать формат даты '%s'\n", *dateStr)
			return
		}
	} else {
		now = time.Now()
	}

	// Конвертация в UTC если нужно
	if *utc {
		now = now.UTC()
	}

	// Вывод в выбранном формате
	switch *format {
	case "rfc3339":
		fmt.Println(now.Format(time.RFC3339))
	case "unix":
		fmt.Println(now.Unix())
	case "kitchen":
		fmt.Println(now.Format(time.Kitchen))
	case "date":
		fmt.Println(now.Format("2006-01-02"))
	case "time":
		fmt.Println(now.Format("15:04:05"))
	default:
		fmt.Println(now.Format("2006-01-02 15:04:05 Monday"))
	}
}