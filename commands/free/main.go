package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syscall"
)

// MemInfo хранит информацию о памяти
type MemInfo struct {
	Total     uint64
	Free      uint64
	Available uint64
	Used      uint64
	SwapTotal uint64
	SwapFree  uint64
}

func main() {
	// Определение флагов
	var (
		help   = flag.Bool("h", false, "показать справку")
		human  = flag.Bool("human", false, "человеко-читаемый вывод")
		bytes  = flag.Bool("b", false, "вывод в байтах")
		swap   = flag.Bool("s", false, "показать swap")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("free - информация о памяти")
		fmt.Println("Использование: free [-human] [-b] [-s]")
		fmt.Println("  -human   человеко-читаемый формат")
		fmt.Println("  -b       вывод в байтах")
		fmt.Println("  -s       показать swap")
		return
	}

	// Получаем информацию о памяти
	memInfo, err := getMemoryInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "free: %v\n", err)
		os.Exit(1)
	}

	// Выводим информацию
	if *human {
		printHumanReadable(memInfo, *swap)
	} else if *bytes {
		printBytes(memInfo, *swap)
	} else {
		printDefault(memInfo, *swap)
	}
}

// getMemoryInfo получает информацию о памяти через syscall
func getMemoryInfo() (*MemInfo, error) {
	var sysInfo syscall.Sysinfo_t
	if err := syscall.Sysinfo(&sysInfo); err != nil {
		return nil, err
	}

	// Конвертируем в байты
	unit := uint64(sysInfo.Unit)
	totalRAM := uint64(sysInfo.Totalram) * unit
	freeRAM := uint64(sysInfo.Freeram) * unit
	sharedRAM := uint64(sysInfo.Sharedram) * unit
	bufferRAM := uint64(sysInfo.Bufferram) * unit
	
	// Примерное вычисление available (в Linux обычно больше)
	available := freeRAM + bufferRAM
	
	// Для swap
	totalSwap := uint64(sysInfo.Totalswap) * unit
	freeSwap := uint64(sysInfo.Freeswap) * unit

	return &MemInfo{
		Total:     totalRAM,
		Free:      freeRAM,
		Available: available,
		Used:      totalRAM - available,
		SwapTotal: totalSwap,
		SwapFree:  freeSwap,
	}, nil
}

// printDefault выводит информацию в стандартном формате
func printDefault(mem *MemInfo, showSwap bool) {
	fmt.Printf("              total        used        free      available\n")
	fmt.Printf("Mem:     %10d  %10d  %10d  %10d\n", 
		mem.Total/1024, mem.Used/1024, mem.Free/1024, mem.Available/1024)
	
	if showSwap {
		swapUsed := mem.SwapTotal - mem.SwapFree
		fmt.Printf("Swap:    %10d  %10d  %10d\n", 
			mem.SwapTotal/1024, swapUsed/1024, mem.SwapFree/1024)
	}
}

// printBytes выводит информацию в байтах
func printBytes(mem *MemInfo, showSwap bool) {
	fmt.Printf("Mem: Total=%d Used=%d Free=%d Available=%d\n", 
		mem.Total, mem.Used, mem.Free, mem.Available)
	
	if showSwap {
		swapUsed := mem.SwapTotal - mem.SwapFree
		fmt.Printf("Swap: Total=%d Used=%d Free=%d\n", 
			mem.SwapTotal, swapUsed, mem.SwapFree)
	}
}

// printHumanReadable выводит информацию в человеко-читаемом формате
func printHumanReadable(mem *MemInfo, showSwap bool) {
	fmt.Printf("Mem:     Total=%s  Used=%s  Free=%s  Available=%s\n", 
		formatBytes(mem.Total), formatBytes(mem.Used), 
		formatBytes(mem.Free), formatBytes(mem.Available))
	
	if showSwap {
		swapUsed := mem.SwapTotal - mem.SwapFree
		fmt.Printf("Swap:    Total=%s  Used=%s  Free=%s\n", 
			formatBytes(mem.SwapTotal), formatBytes(swapUsed), formatBytes(mem.SwapFree))
	}
	
	// Дополнительная информация из runtime
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Go Heap: Alloc=%s  Sys=%s  GC=%d\n", 
		formatBytes(m.Alloc), formatBytes(m.Sys), m.NumGC)
}

// formatBytes форматирует байты в человеко-читаемый вид
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