package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type MemInfo struct {
	Total     uint64
	Free      uint64
	Available uint64
	Used      uint64
	SwapTotal uint64
	SwapFree  uint64
	SwapUsed  uint64
	Buffers   uint64
	Cached    uint64
}

func main() {
	var (
		help     = flag.Bool("h", false, "показать справку")
		human    = flag.Bool("h", false, "человеко-читаемый вывод")
		bytes    = flag.Bool("b", false, "вывод в байтах")
		mega     = flag.Bool("m", false, "вывод в мегабайтах")
		giga     = flag.Bool("g", false, "вывод в гигабайтах")
		showSwap = flag.Bool("s", false, "показать swap")
	)
	flag.Parse()

	if *help {
		fmt.Println("free - информация о памяти")
		fmt.Println("Использование: free [-b] [-m] [-g] [-h] [-s]")
		fmt.Println("  -b    вывод в байтах")
		fmt.Println("  -m    вывод в мегабайтах")
		fmt.Println("  -g    вывод в гигабайтах")
		fmt.Println("  -h    человеко-читаемый вывод")
		fmt.Println("  -s    показать swap")
		return
	}

	memInfo, err := getMemoryInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "free: %v\n", err)
		os.Exit(1)
	}

	// Определяем формат вывода
	if *bytes {
		printInBytes(memInfo, *showSwap)
	} else if *mega {
		printInMB(memInfo, *showSwap)
	} else if *giga {
		printInGB(memInfo, *showSwap)
	} else if *human {
		printHumanReadable(memInfo, *showSwap)
	} else {
		printDefault(memInfo, *showSwap)
	}
}

// getMemoryInfo получает информацию о памяти через /proc/meminfo (более точно)
func getMemoryInfo() (*MemInfo, error) {
	// Читаем /proc/meminfo для более точных данных (Linux)
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		// Fallback на syscall если /proc/meminfo недоступен
		return getMemoryInfoSyscall()
	}

	return parseMemInfo(string(data))
}

// getMemoryInfoSyscall получает информацию через syscall (запасной вариант)
func getMemoryInfoSyscall() (*MemInfo, error) {
	var sysInfo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&sysInfo); err != nil {
		return nil, err
	}

	unit := uint64(sysInfo.Unit)
	totalRAM := uint64(sysInfo.Totalram) * unit
	freeRAM := uint64(sysInfo.Freeram) * unit
	sharedRAM := uint64(sysInfo.Sharedram) * unit
	bufferRAM := uint64(sysInfo.Bufferram) * unit

	// Приблизительный расчет
	available := freeRAM + bufferRAM

	return &MemInfo{
		Total:     totalRAM,
		Free:      freeRAM,
		Available: available,
		Used:      totalRAM - available,
		SwapTotal: uint64(sysInfo.Totalswap) * unit,
		SwapFree:  uint64(sysInfo.Freeswap) * unit,
		SwapUsed:  (uint64(sysInfo.Totalswap) - uint64(sysInfo.Freeswap)) * unit,
	}, nil
}

// parseMemInfo парсит /proc/meminfo
func parseMemInfo(content string) (*MemInfo, error) {
	info := &MemInfo{}
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := strings.TrimSuffix(fields[0], ":")
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		// Конвертируем из KB в байты
		value *= 1024

		switch key {
		case "MemTotal":
			info.Total = value
		case "MemFree":
			info.Free = value
		case "MemAvailable":
			info.Available = value
		case "Buffers":
			info.Buffers = value
		case "Cached":
			info.Cached = value
		case "SwapTotal":
			info.SwapTotal = value
		case "SwapFree":
			info.SwapFree = value
		}
	}

	// Вычисляем Used
	if info.Available > 0 {
		info.Used = info.Total - info.Available
	} else {
		info.Used = info.Total - info.Free
	}

	info.SwapUsed = info.SwapTotal - info.SwapFree

	return info, nil
}

// printDefault выводит в стандартном формате (в KB)
func printDefault(mem *MemInfo, showSwap bool) {
	fmt.Printf("%-10s %12s %12s %12s %12s\n", "", "total", "used", "free", "available")
	fmt.Printf("%-10s %12d %12d %12d %12d\n",
		"Mem:",
		mem.Total/1024,
		mem.Used/1024,
		mem.Free/1024,
		mem.Available/1024)

	if showSwap {
		fmt.Printf("%-10s %12d %12d %12d\n",
			"Swap:",
			mem.SwapTotal/1024,
			mem.SwapUsed/1024,
			mem.SwapFree/1024)
	}
}

// printInBytes выводит в байтах
func printInBytes(mem *MemInfo, showSwap bool) {
	fmt.Printf("%-10s %12d %12d %12d %12d\n", "", "total", "used", "free", "available")
	fmt.Printf("%-10s %12d %12d %12d %12d\n",
		"Mem:", mem.Total, mem.Used, mem.Free, mem.Available)

	if showSwap {
		fmt.Printf("%-10s %12d %12d %12d\n", "Swap:", mem.SwapTotal, mem.SwapUsed, mem.SwapFree)
	}
}

// printInMB выводит в мегабайтах
func printInMB(mem *MemInfo, showSwap bool) {
	toMB := func(bytes uint64) uint64 {
		return bytes / (1024 * 1024)
	}

	fmt.Printf("%-10s %12s %12s %12s %12s\n", "", "total", "used", "free", "available")
	fmt.Printf("%-10s %12d %12d %12d %12d\n",
		"Mem:",
		toMB(mem.Total),
		toMB(mem.Used),
		toMB(mem.Free),
		toMB(mem.Available))

	if showSwap {
		fmt.Printf("%-10s %12d %12d %12d\n", "Swap:", toMB(mem.SwapTotal), toMB(mem.SwapUsed), toMB(mem.SwapFree))
	}
}

// printInGB выводит в гигабайтах
func printInGB(mem *MemInfo, showSwap bool) {
	toGB := func(bytes uint64) uint64 {
		return bytes / (1024 * 1024 * 1024)
	}

	fmt.Printf("%-10s %12s %12s %12s %12s\n", "", "total", "used", "free", "available")
	fmt.Printf("%-10s %12d %12d %12d %12d\n",
		"Mem:",
		toGB(mem.Total),
		toGB(mem.Used),
		toGB(mem.Free),
		toGB(mem.Available))

	if showSwap {
		fmt.Printf("%-10s %12d %12d %12d\n", "Swap:", toGB(mem.SwapTotal), toGB(mem.SwapUsed), toGB(mem.SwapFree))
	}
}

// printHumanReadable выводит в человеко-читаемом формате
func printHumanReadable(mem *MemInfo, showSwap bool) {
	fmt.Printf("%-10s %12s %12s %12s %12s\n", "", "total", "used", "free", "available")
	fmt.Printf("%-10s %12s %12s %12s %12s\n",
		"Mem:",
		formatBytes(mem.Total),
		formatBytes(mem.Used),
		formatBytes(mem.Free),
		formatBytes(mem.Available))

	if showSwap {
		fmt.Printf("%-10s %12s %12s %12s\n", "Swap:", formatBytes(mem.SwapTotal), formatBytes(mem.SwapUsed), formatBytes(mem.SwapFree))
	}

	// Дополнительная информация о Go runtime
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nGo Runtime:\n")
	fmt.Printf("  Heap Alloc: %s\n", formatBytes(m.Alloc))
	fmt.Printf("  Heap Sys:   %s\n", formatBytes(m.HeapSys))
	fmt.Printf("  GC Count:   %d\n", m.NumGC)
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"K", "M", "G", "T", "P", "E"}
	return fmt.Sprintf("%.1f %sB", float64(bytes)/float64(div), units[exp])
}
