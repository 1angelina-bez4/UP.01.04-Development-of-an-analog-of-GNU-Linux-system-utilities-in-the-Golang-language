package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	PID     int
	PPID    int
	Name    string
	State   string
	CPU     float64
	Memory  float64
	User    string
}

func main() {
	// Определение флагов
	var (
		help     = flag.Bool("h", false, "показать справку")
		all      = flag.Bool("e", false, "показать все процессы")
		user     = flag.String("u", "", "показать процессы пользователя")
		pid      = flag.Int("p", 0, "показать процесс с указанным PID")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("ps - информация о процессах")
		fmt.Println("Использование: ps [-e] [-u пользователь] [-p PID]")
		fmt.Println("  -e    показать все процессы")
		fmt.Println("  -u    показать процессы пользователя")
		fmt.Println("  -p    показать процесс с указанным PID")
		return
	}

	// Получаем список процессов
	processes, err := getProcesses()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ps: %v\n", err)
		os.Exit(1)
	}

	// Фильтрация процессов
	var filtered []Process
	for _, proc := range processes {
		if *pid > 0 && proc.PID != *pid {
			continue
		}
		if *user != "" && proc.User != *user {
			continue
		}
		if !*all && proc.PID == 1 {
			continue
		}
		filtered = append(filtered, proc)
	}

	// Вывод заголовка
	fmt.Printf("%-8s %-8s %-8s %-10s %-8s %s\n", 
		"PID", "PPID", "STATE", "CPU%", "MEM%", "COMMAND")
	fmt.Println(strings.Repeat("-", 60))

	// Вывод процессов
	for _, proc := range filtered {
		fmt.Printf("%-8d %-8d %-8s %-10.1f %-8.1f %s\n",
			proc.PID, proc.PPID, proc.State, proc.CPU, proc.Memory, proc.Name)
	}
}

// getProcesses возвращает список процессов в Linux
func getProcesses() ([]Process, error) {
	var processes []Process
	
	// Читаем директорию /proc
	entries, err := ioutil.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		// Проверяем, является ли имя директории числом (PID)
		pid, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		
		// Читаем информацию о процессе
		proc := Process{PID: pid}
		
		// Читаем статус процесса
		statusFile := fmt.Sprintf("/proc/%d/status", pid)
		statusData, err := ioutil.ReadFile(statusFile)
		if err == nil {
			lines := strings.Split(string(statusData), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PPid:") {
					parts := strings.Fields(line)
					if len(parts) > 1 {
						proc.PPID, _ = strconv.Atoi(parts[1])
					}
				} else if strings.HasPrefix(line, "State:") {
					parts := strings.Fields(line)
					if len(parts) > 1 {
						proc.State = parts[1]
					}
				} else if strings.HasPrefix(line, "Name:") {
					parts := strings.Fields(line)
					if len(parts) > 1 {
						proc.Name = parts[1]
					}
				} else if strings.HasPrefix(line, "Uid:") {
					parts := strings.Fields(line)
					if len(parts) > 1 {
						// Упрощённо - в реальности нужно преобразовывать UID в имя
						proc.User = parts[1]
					}
				}
			}
		}
		
		// Читаем статистику CPU и памяти
		statFile := fmt.Sprintf("/proc/%d/stat", pid)
		statData, err := ioutil.ReadFile(statFile)
		if err == nil {
			parts := strings.Fields(string(statData))
			if len(parts) > 22 {
				// Упрощённый расчёт CPU
				proc.CPU = 0.0
				// RSS в страницах (обычно 4KB на страницу)
				rss, _ := strconv.ParseFloat(parts[23], 64)
				proc.Memory = rss * 4 / 1024 // Примерно в MB
			}
		}
		
		if proc.Name == "" {
			proc.Name = fmt.Sprintf("[%d]", pid)
		}
		
		processes = append(processes, proc)
	}
	
	return processes, nil
}