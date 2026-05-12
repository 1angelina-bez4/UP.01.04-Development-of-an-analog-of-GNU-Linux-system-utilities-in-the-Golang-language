package main

import (
	"flag"
	"fmt"
	"runtime"
	"os"
)

type ArchInfo struct {
	Arch      string
	OS        string
	CPUCores  int
	GoVersion string
}

func main() {
	// Определение флагов
	var (
		help     = flag.Bool("h", false, "показать справку")
		long     = flag.Bool("l", false, "подробный вывод")
		json     = flag.Bool("json", false, "вывод в формате JSON")
		all      = flag.Bool("a", false, "показать всю информацию")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("arch - вывод архитектуры системы")
		fmt.Println("Использование: arch [-l] [-json] [-a]")
		fmt.Println("  -l     подробный вывод")
		fmt.Println("  -json  вывод в формате JSON")
		fmt.Println("  -a     показать всю информацию")
		return
	}

	info := ArchInfo{
		Arch:      runtime.GOARCH,
		OS:        runtime.GOOS,
		CPUCores:  runtime.NumCPU(),
		GoVersion: runtime.Version(),
	}

	// Вывод в зависимости от флагов
	if *json {
		printJSON(info)
	} else if *long || *all {
		printLong(info)
	} else {
		fmt.Println(info.Arch)
	}
}

func printLong(info ArchInfo) {
	fmt.Printf("Архитектура:    %s\n", info.Arch)
	fmt.Printf("Операционная система: %s\n", info.OS)
	fmt.Printf("Количество CPU: %d\n", info.CPUCores)
	fmt.Printf("Версия Go:      %s\n", info.GoVersion)
}

func printJSON(info ArchInfo) {
	fmt.Printf("{\n")
	fmt.Printf("  \"architecture\": \"%s\",\n", info.Arch)
	fmt.Printf("  \"os\": \"%s\",\n", info.OS)
	fmt.Printf("  \"cpu_cores\": %d,\n", info.CPUCores)
	fmt.Printf("  \"go_version\": \"%s\"\n", info.GoVersion)
	fmt.Printf("}\n")
}