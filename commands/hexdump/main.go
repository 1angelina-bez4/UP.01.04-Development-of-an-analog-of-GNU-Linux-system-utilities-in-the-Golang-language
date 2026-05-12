package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type HexOptions struct {
	Canonical bool
	Bytes     int
	Offset    int64
	Length    int64
}

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		canonical = flag.Bool("C", false, "канонический формат (hex+ascii)")
		bytes     = flag.Int("n", 16, "байт на строку")
		offset    = flag.Int64("s", 0, "начальное смещение")
		length    = flag.Int64("L", -1, "количество байт для вывода")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("hexdump - шестнадцатеричный дамп файла")
		fmt.Println("Использование: hexdump [-C] [-n байт] [-s смещение] [-L длина] [файл]")
		fmt.Println("  -C    канонический формат (hex+ascii)")
		fmt.Println("  -n    количество байт на строку (по умолчанию 16)")
		fmt.Println("  -s    начальное смещение")
		fmt.Println("  -L    общее количество байт для вывода")
		return
	}

	if *bytes < 1 || *bytes > 64 {
		fmt.Fprintln(os.Stderr, "hexdump: количество байт на строку должно быть от 1 до 64")
		os.Exit(1)
	}

	var file *os.File
	var err error

	if flag.NArg() > 0 {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "hexdump: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
	} else {
		file = os.Stdin
	}

	opts := HexOptions{
		Canonical: *canonical,
		Bytes:     *bytes,
		Offset:    *offset,
		Length:    *length,
	}

	dumpHex(file, opts)
}

// dumpHex выводит шестнадцатеричный дамп
func dumpHex(file *io.File, opts HexOptions) {
	// Устанавливаем смещение
	if opts.Offset > 0 {
		file.Seek(opts.Offset, io.SeekStart)
	}

	buffer := make([]byte, opts.Bytes)
	position := opts.Offset
	totalRead := int64(0)

	for {
		// Определяем, сколько читать
		toRead := opts.Bytes
		if opts.Length > 0 {
			remaining := opts.Length - totalRead
			if remaining <= 0 {
				break
			}
			if int(remaining) < toRead {
				toRead = int(remaining)
			}
		}

		n, err := file.Read(buffer[:toRead])
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "hexdump: ошибка чтения: %v\n", err)
			break
		}

		if opts.Canonical {
			printCanonical(buffer[:n], position)
		} else {
			printSimple(buffer[:n], position)
		}

		position += int64(n)
		totalRead += int64(n)
	}
}

// printCanonical выводит в каноническом формате (hex + ascii)
func printCanonical(data []byte, offset int64) {
	// Вывод смещения
	fmt.Printf("%08x  ", offset)

	// Вывод шестнадцатеричных байт
	for i, b := range data {
		fmt.Printf("%02x", b)
		if (i+1)%8 == 0 {
			fmt.Print(" ")
		}
		if (i+1)%16 == 0 && i+1 < len(data) {
			fmt.Print(" ")
		}
	}

	// Выравнивание
	bytesPerLine := 16
	if len(data) < bytesPerLine {
		spaces := (bytesPerLine - len(data)) * 3
		if (bytesPerLine - len(data)) > 8 {
			spaces += 1
		}
		for i := 0; i < spaces; i++ {
			fmt.Print(" ")
		}
	}

	fmt.Print("  |")

	// Вывод ASCII представления
	for _, b := range data {
		if b >= 32 && b <= 126 {
			fmt.Printf("%c", b)
		} else {
			fmt.Print(".")
		}
	}

	fmt.Println("|")
}

// printSimple выводит в простом формате
func printSimple(data []byte, offset int64) {
	fmt.Printf("%08x: ", offset)
	for i, b := range data {
		fmt.Printf("%02x", b)
		if (i+1)%8 == 0 && i+1 < len(data) {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}