package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const historyFile = ".goutils_history"

type HistoryEntry struct {
	Number int
	Command string
}

func main() {
	// Определение флагов
	var (
		help     = flag.Bool("h", false, "показать справку")
		clear    = flag.Bool("c", false, "очистить историю")
		delete   = flag.Int("d", 0, "удалить команду по номеру")
		search   = flag.String("s", "", "поиск команд по подстроке")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("history - отображение истории команд")
		fmt.Println("Использование: history [-c] [-d номер] [-s строка]")
		fmt.Println("  -c    очистить историю")
		fmt.Println("  -d    удалить команду по номеру")
		fmt.Println("  -s    показать только команды, содержащие строку")
		return
	}

	historyPath := getHistoryPath()

	// Очистка истории
	if *clear {
		if err := os.Remove(historyPath); err != nil {
			fmt.Fprintf(os.Stderr, "history: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("history: история очищена")
		return
	}

	// Удаление конкретной команды
	if *delete > 0 {
		if err := deleteCommand(historyPath, *delete); err != nil {
			fmt.Fprintf(os.Stderr, "history: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("history: удалена команда #%d\n", *delete)
		return
	}

	// Чтение и отображение истории
	history, err := readHistory(historyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "history: %v\n", err)
		os.Exit(1)
	}

	// Поиск по подстроке
	for _, entry := range history {
		if *search == "" || strings.Contains(entry.Command, *search) {
			fmt.Printf("%5d  %s\n", entry.Number, entry.Command)
		}
	}
}

// getHistoryPath возвращает путь к файлу истории
func getHistoryPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return historyFile
	}
	return filepath.Join(home, historyFile)
}

// readHistory читает историю из файла
func readHistory(path string) ([]HistoryEntry, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var history []HistoryEntry
	scanner := bufio.NewScanner(file)
	num := 1

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			history = append(history, HistoryEntry{
				Number:  num,
				Command: line,
			})
			num++
		}
	}

	return history, scanner.Err()
}

// deleteCommand удаляет команду по номеру
func deleteCommand(path string, num int) error {
	history, err := readHistory(path)
	if err != nil {
		return err
	}

	if num < 1 || num > len(history) {
		return fmt.Errorf("неверный номер команды: %d", num)
	}

	// Создаём новую историю без удалённой команды
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range history {
		if entry.Number != num {
			fmt.Fprintln(writer, entry.Command)
		}
	}
	writer.Flush()

	return nil
}

// AddCommand добавляет команду в историю
func AddCommand(cmd string) {
	path := getHistoryPath()
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	
	fmt.Fprintln(file, cmd)
}