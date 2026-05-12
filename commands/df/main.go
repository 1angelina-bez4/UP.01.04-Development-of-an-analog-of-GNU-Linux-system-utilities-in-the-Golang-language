package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

type Filesystem struct {
	Device     string
	MountPoint string
	Total      uint64
	Used       uint64
	Free       uint64
	UsePercent float64
}

func main() {
	// Определение флагов
	var (
		help    = flag.Bool("h", false, "показать справку")
		human   = flag.Bool("human", false, "человеко-читаемый вывод")
		total   = flag.Bool("total", false, "показать общий итог")
		local   = flag.Bool("l", false, "только локальные файловые системы")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("df - отчёт об использовании дискового пространства")
		fmt.Println("Использование: df [-human] [-total] [-l] [путь...]")
		fmt.Println("  -human    человеко-читаемый формат")
		fmt.Println("  -total    показать общий итог")
		fmt.Println("  -l        только локальные файловые системы")
		return
	}

	// Получаем информацию о файловых системах
	filesystems := getFilesystems(*local)

	// Фильтруем по указанным путям
	if flag.NArg() > 0 {
		var filtered []Filesystem
		paths := make(map[string]bool)
		for _, p := range flag.Args() {
			paths[p] = true
		}
		for _, fs := range filesystems {
			if paths[fs.MountPoint] {
				filtered = append(filtered, fs)
			}
		}
		filesystems = filtered
	}

	// Выводим результат
	printFilesystems(filesystems, *human)

	// Общий итог
	if *total && len(filesystems) > 1 {
		var totalTotal, totalUsed, totalFree uint64
		for _, fs := range filesystems {
			totalTotal += fs.Total
			totalUsed += fs.Used
			totalFree += fs.Free
		}
		fmt.Println("\nTotal:")
		if *human {
			fmt.Printf("  Total: %s\n", formatBytes(totalTotal))
			fmt.Printf("  Used:  %s\n", formatBytes(totalUsed))
			fmt.Printf("  Free:  %s\n", formatBytes(totalFree))
		} else {
			fmt.Printf("  Total: %d\n", totalTotal)
			fmt.Printf("  Used:  %d\n", totalUsed)
			fmt.Printf("  Free:  %d\n", totalFree)
		}
	}
}

// getFilesystems возвращает список файловых систем
func getFilesystems(localOnly bool) []Filesystem {
	var filesystems []Filesystem
	
	// В Linux мы можем прочитать /proc/mounts
	mountsFile, err := os.Open("/proc/mounts")
	if err != nil {
		fmt.Fprintf(os.Stderr, "df: %v\n", err)
		return filesystems
	}
	defer mountsFile.Close()

	// Читаем информацию о каждой точке монтирования
	var stat syscall.Statfs_t
	
	// Упрощённая версия - реальные точки монтирования
	mountPoints := []string{"/", "/home", "/boot", "/var", "/tmp"}
	
	for _, mp := range mountPoints {
		err := syscall.Statfs(mp, &stat)
		if err != nil {
			continue
		}
		
		total := stat.Blocks * uint64(stat.Bsize)
		free := stat.Bfree * uint64(stat.Bsize)
		available := stat.Bavail * uint64(stat.Bsize)
		used := total - free
		usePercent := float64(used) / float64(total) * 100
		
		filesystems = append(filesystems, Filesystem{
			Device:     mp,
			MountPoint: mp,
			Total:      total,
			Used:       used,
			Free:       available,
			UsePercent: usePercent,
		})
	}
	
	return filesystems
}

// printFilesystems выводит информацию о файловых системах
func printFilesystems(fsList []Filesystem, human bool) {
	if human {
		fmt.Printf("%-20s %10s %10s %10s %8s %s\n", 
			"Filesystem", "Size", "Used", "Avail", "Use%", "Mounted on")
	} else {
		fmt.Printf("%-20s %12s %12s %12s %8s %s\n", 
			"Filesystem", "1K-blocks", "Used", "Available", "Use%", "Mounted on")
	}
	
	for _, fs := range fsList {
		if human {
			fmt.Printf("%-20s %10s %10s %10s %7.1f%% %s\n",
				fs.Device,
				formatBytes(fs.Total),
				formatBytes(fs.Used),
				formatBytes(fs.Free),
				fs.UsePercent,
				fs.MountPoint)
		} else {
			fmt.Printf("%-20s %12d %12d %12d %7.1f%% %s\n",
				fs.Device,
				fs.Total/1024,
				fs.Used/1024,
				fs.Free/1024,
				fs.UsePercent,
				fs.MountPoint)
		}
	}
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
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}