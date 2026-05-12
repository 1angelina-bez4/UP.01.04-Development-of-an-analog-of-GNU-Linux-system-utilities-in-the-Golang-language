package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

type Signal struct {
	Name string
	Num  int
}

var signals = map[string]int{
	"HUP":  1,
	"INT":  2,
	"QUIT": 3,
	"KILL": 9,
	"TERM": 15,
}

func main() {
	// Определение флагов
	var (
		help     = flag.Bool("h", false, "показать справку")
		signal   = flag.String("s", "TERM", "сигнал для отправки")
		list     = flag.Bool("l", false, "список сигналов")
		verbose  = flag.Bool("v", false, "подробный вывод")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("kill - отправка сигналов процессам")
		fmt.Println("Использование: kill [-s сигнал] [-l] [-v] PID...")
		fmt.Println("  -s    сигнал (TERM, KILL, HUP, INT, QUIT)")
		fmt.Println("  -l    показать список сигналов")
		fmt.Println("  -v    подробный вывод")
		return
	}

	// Показываем список сигналов
	if *list {
		fmt.Println("Доступные сигналы:")
		for name, num := range signals {
			fmt.Printf("  %s (%d)\n", name, num)
		}
		return
	}

	// Проверка наличия PID
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "kill: не указан PID")
		os.Exit(1)
	}

	// Получаем номер сигнала
	sigNum, ok := signals[*signal]
	if !ok {
		fmt.Fprintf(os.Stderr, "kill: неизвестный сигнал '%s'\n", *signal)
		os.Exit(1)
	}

	// Отправляем сигнал каждому процессу
	for _, pidStr := range flag.Args() {
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: неверный PID '%s'\n", pidStr)
			continue
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: %v\n", err)
			continue
		}

		err = process.Signal(syscall.Signal(sigNum))
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: не удалось отправить сигнал %d процессу %d: %v\n", 
				sigNum, pid, err)
		} else if *verbose {
			fmt.Printf("kill: сигнал %s (%d) отправлен процессу %d\n", *signal, sigNum, pid)
		}
	}
}