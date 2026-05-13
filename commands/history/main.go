package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const histFile = ".goutils_history"

func main() {
	add := flag.String("add", "", "добавить команду")
	help := flag.Bool("h", false, "справка")
	flag.Parse()

	if *help {
		fmt.Println("history - управление историей")
		fmt.Println("  history               - показать историю")
		fmt.Println("  history -n N          - последние N команд")
		fmt.Println("  history -add 'cmd'    - добавить команду")
		fmt.Println("  history -x N          - выполнить команду #N")
		fmt.Println("  history -d N          - удалить команду #N")
		fmt.Println("  history -c            - очистить историю")
		return
	}

	// Если передан флаг -add
	if *add != "" {
		AddCommand(*add)
		fmt.Printf("✓ Добавлено: %s\n", *add)
		return
	}

	// Если есть аргумент без флага, считаем его командой для выполнения
	if flag.NArg() > 0 && flag.Arg(0) != "" {
		cmd := strings.Join(flag.Args(), " ")
		AddCommand(cmd)

		// Выполняем команду
		parts := strings.Fields(cmd)
		if len(parts) > 0 {
			c := exec.Command(parts[0], parts[1:]...)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Run()
		}
		return
	}

	// Обычный вывод истории
	var (
		clear   = flag.Bool("c", false, "")
		delete  = flag.Int("d", 0, "")
		execNum = flag.Int("x", 0, "")
		last    = flag.Int("n", 0, "")
	)
	flag.Parse()

	path := getPath()

	switch {
	case *clear:
		os.Remove(path)
		fmt.Println("history: очищена")
		return
	case *delete > 0:
		deleteCmd(path, *delete)
		return
	case *execNum > 0:
		execCmdNum(path, *execNum)
		return
	}

	history := readHistory(path)
	if *last > 0 && *last < len(history) {
		history = history[len(history)-*last:]
	}

	if len(history) == 0 {
		fmt.Println("История пуста. Добавьте команды через:")
		fmt.Println("  history -add 'команда'")
		fmt.Println("  или")
		fmt.Println("  history 'команда'")
		return
	}

	for _, h := range history {
		fmt.Printf("%6d  %s\n", h.num, h.cmd)
	}
}

type entry struct {
	num int
	cmd string
}

func getPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, histFile)
}

func readHistory(path string) []entry {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	var history []entry
	scanner := bufio.NewScanner(file)
	num := 1
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			history = append(history, entry{num, line})
			num++
		}
	}
	return history
}

func writeHistory(path string, history []entry) {
	file, _ := os.Create(path)
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, h := range history {
		fmt.Fprintln(w, h.cmd)
	}
	w.Flush()
}

func deleteCmd(path string, num int) {
	history := readHistory(path)
	if num < 1 || num > len(history) {
		fmt.Fprintf(os.Stderr, "history: неверный номер %d\n", num)
		return
	}
	newHistory := []entry{}
	for _, h := range history {
		if h.num != num {
			newHistory = append(newHistory, h)
		}
	}
	for i := range newHistory {
		newHistory[i].num = i + 1
	}
	writeHistory(path, newHistory)
	fmt.Printf("history: удалена #%d\n", num)
}

func execCmdNum(path string, num int) {
	history := readHistory(path)
	if num < 1 || num > len(history) {
		fmt.Fprintf(os.Stderr, "history: неверный номер %d\n", num)
		return
	}
	cmd := history[num-1].cmd
	fmt.Printf("Выполнение: %s\n", cmd)

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	c := exec.Command(parts[0], parts[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	c.Run()
}

// AddCommand добавляет команду в историю
func AddCommand(cmd string) {
	if cmd == "" {
		return
	}
	path := getPath()
	history := readHistory(path)

	history = append(history, entry{len(history) + 1, cmd})

	// Оставляем последние 1000
	if len(history) > 1000 {
		history = history[len(history)-1000:]
	}

	// Перенумерация
	for i := range history {
		history[i].num = i + 1
	}

	writeHistory(path, history)
}
