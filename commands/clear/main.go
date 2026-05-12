package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	// Определение флагов
	var (
		help     = flag.Bool("h", false, "показать справку")
		keepLine = flag.Bool("k", false, "сохранить строку с приглашением")
		scrollback = flag.Bool("s", false, "очистить буфер прокрутки")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("clear - очистка экрана терминала")
		fmt.Println("Использование: clear [-k] [-s]")
		fmt.Println("  -k    сохранить строку с приглашением")
		fmt.Println("  -s    очистить буфер прокрутки")
		return
	}

	// Очистка экрана
	if *scrollback {
		// Очистка с буфером прокрутки
		switch runtime.GOOS {
		case "windows":
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		default:
			fmt.Print("\033[3J\033[H\033[2J")
		}
	} else {
		// Обычная очистка экрана
		fmt.Print("\033[2J\033[H")
	}

	// Если нужно сохранить строку
	if *keepLine {
		fmt.Print("\033[7m> \033[0m")
	}
}