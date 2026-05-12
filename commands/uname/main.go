package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syscall"
)

type SystemInfo struct {
	KernelName    string
	NodeName      string
	KernelRelease string
	KernelVersion string
	Machine       string
	OperatingSystem string
}

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		all       = flag.Bool("a", false, "показать всё")
		kernel    = flag.Bool("s", false, "имя ядра")
		node      = flag.Bool("n", false, "имя узла")
		release   = flag.Bool("r", false, "версия ядра")
		version   = flag.Bool("v", false, "версия ядра (подробно)")
		machine   = flag.Bool("m", false, "архитектура")
		osname    = flag.Bool("o", false, "ОС")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("uname - вывод информации о системе")
		fmt.Println("Использование: uname [-a] [-s] [-n] [-r] [-v] [-m] [-o]")
		fmt.Println("  -a    всё")
		fmt.Println("  -s    имя ядра")
		fmt.Println("  -n    имя узла")
		fmt.Println("  -r    релиз ядра")
		fmt.Println("  -v    версия ядра")
		fmt.Println("  -m    архитектура")
		fmt.Println("  -o    ОС")
		return
	}

	info := getSystemInfo()

	// Если указан -a или нет других флагов, показываем всё
	if *all || (!*kernel && !*node && !*release && !*version && !*machine && !*osname) {
		fmt.Printf("%s %s %s %s %s %s\n",
			info.KernelName,
			info.NodeName,
			info.KernelRelease,
			info.KernelVersion,
			info.Machine,
			info.OperatingSystem)
		return
	}

	// Вывод отдельных частей
	if *kernel {
		fmt.Print(info.KernelName)
	}
	if *node {
		if *kernel { fmt.Print(" ") }
		fmt.Print(info.NodeName)
	}
	if *release {
		if *kernel || *node { fmt.Print(" ") }
		fmt.Print(info.KernelRelease)
	}
	if *version {
		if *kernel || *node || *release { fmt.Print(" ") }
		fmt.Print(info.KernelVersion)
	}
	if *machine {
		if *kernel || *node || *release || *version { fmt.Print(" ") }
		fmt.Print(info.Machine)
	}
	if *osname {
		if *kernel || *node || *release || *version || *machine { fmt.Print(" ") }
		fmt.Print(info.OperatingSystem)
	}
	fmt.Println()
}

// getSystemInfo собирает информацию о системе
func getSystemInfo() SystemInfo {
	var uts syscall.Utsname
	if err := syscall.Uname(&uts); err != nil {
		fmt.Fprintf(os.Stderr, "uname: %v\n", err)
		os.Exit(1)
	}

	return SystemInfo{
		KernelName:    charsToString(uts.Sysname[:]),
		NodeName:      charsToString(uts.Nodename[:]),
		KernelRelease: charsToString(uts.Release[:]),
		KernelVersion: charsToString(uts.Version[:]),
		Machine:       charsToString(uts.Machine[:]),
		OperatingSystem: runtime.GOOS,
	}
}

// charsToString преобразует массив байт в строку
func charsToString(ca []int8) string {
	s := make([]byte, len(ca))
	len := 0
	for _, v := range ca {
		if v == 0 {
			break
		}
		s[len] = byte(v)
		len++
	}
	return string(s[:len])
}