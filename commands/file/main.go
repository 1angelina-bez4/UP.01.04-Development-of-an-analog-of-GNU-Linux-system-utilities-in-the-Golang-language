package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
)

// FileType представляет информацию о типе файла
type FileType struct {
	Name        string
	MimeType    string
	Description string
}

func main() {
	// Определение флагов
	var (
		help      = flag.Bool("h", false, "показать справку")
		mime      = flag.Bool("i", false, "вывести MIME-тип")
		verbose   = flag.Bool("v", false, "подробный вывод")
		extension = flag.Bool("e", false, "определять по расширению")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("file - определение типа файла")
		fmt.Println("Использование: file [-i] [-v] [-e] файл...")
		fmt.Println("  -i    показать MIME-тип")
		fmt.Println("  -v    подробный вывод")
		fmt.Println("  -e    определить по расширению")
		return
	}

	// Проверка наличия аргументов
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "file: не указан файл")
		os.Exit(1)
	}

	// Анализируем каждый файл
	for _, fname := range flag.Args() {
		info, err := os.Stat(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "file: %v\n", err)
			continue
		}

		var fileType FileType

		if info.IsDir() {
			fileType = FileType{Name: "directory", MimeType: "inode/directory", Description: "директория"}
		} else if *extension {
			fileType = detectByExtension(fname)
		} else {
			fileType = detectByContent(fname)
		}

		// Вывод результата
		if *mime {
			fmt.Printf("%s: %s\n", fname, fileType.MimeType)
		} else if *verbose {
			fmt.Printf("%s: %s (%s)\n", fname, fileType.Description, fileType.MimeType)
		} else {
			fmt.Printf("%s: %s\n", fname, fileType.Description)
		}
	}
}

// detectByContent определяет тип по содержимому файла
func detectByContent(fname string) FileType {
	file, err := os.Open(fname)
	if err != nil {
		return FileType{Name: "unknown", MimeType: "application/octet-stream", Description: "неизвестный"}
	}
	defer file.Close()

	// Читаем первые 512 байт для анализа
	header := make([]byte, 512)
	n, _ := file.Read(header)
	header = header[:n]

	// Проверка на текстовый файл
	if isText(header) {
		return FileType{Name: "text", MimeType: "text/plain", Description: "текстовый файл"}
	}

	// Проверка магических чисел
	magicSignatures := []struct {
		magic  []byte
		offset int
		ftype  FileType
	}{
		{[]byte("\x7fELF"), 0, FileType{Name: "elf", MimeType: "application/x-executable", Description: "ELF исполняемый"}},
		{[]byte("%PDF"), 0, FileType{Name: "pdf", MimeType: "application/pdf", Description: "PDF документ"}},
		{[]byte{0x89, 0x50, 0x4E, 0x47}, 0, FileType{Name: "png", MimeType: "image/png", Description: "PNG изображение"}},
		{[]byte{0xFF, 0xD8}, 0, FileType{Name: "jpg", MimeType: "image/jpeg", Description: "JPEG изображение"}},
		{[]byte("PK\x03\x04"), 0, FileType{Name: "zip", MimeType: "application/zip", Description: "ZIP архив"}},
		{[]byte{0x1F, 0x8B}, 0, FileType{Name: "gzip", MimeType: "application/gzip", Description: "GZIP архив"}},
	}

	for _, sig := range magicSignatures {
		if len(header) >= sig.offset+len(sig.magic) {
			if bytes.Equal(header[sig.offset:sig.offset+len(sig.magic)], sig.magic) {
				return sig.ftype
			}
		}
	}

	return FileType{Name: "data", MimeType: "application/octet-stream", Description: "бинарные данные"}
}

// detectByExtension определяет тип по расширению файла
func detectByExtension(fname string) FileType {
	ext := strings.ToLower(getExtension(fname))

	extMap := map[string]FileType{
		".txt":  {Name: "text", MimeType: "text/plain", Description: "текстовый файл"},
		".go":   {Name: "go", MimeType: "text/x-go", Description: "Go исходник"},
		".py":   {Name: "python", MimeType: "text/x-python", Description: "Python скрипт"},
		".json": {Name: "json", MimeType: "application/json", Description: "JSON данные"},
		".xml":  {Name: "xml", MimeType: "application/xml", Description: "XML документ"},
		".html": {Name: "html", MimeType: "text/html", Description: "HTML документ"},
		".jpg":  {Name: "jpg", MimeType: "image/jpeg", Description: "JPEG изображение"},
		".png":  {Name: "png", MimeType: "image/png", Description: "PNG изображение"},
		".pdf":  {Name: "pdf", MimeType: "application/pdf", Description: "PDF документ"},
	}

	if ft, ok := extMap[ext]; ok {
		return ft
	}
	return FileType{Name: "unknown", MimeType: "application/octet-stream", Description: "неизвестный"}
}

// getExtension возвращает расширение файла
func getExtension(fname string) string {
	for i := len(fname) - 1; i >= 0; i-- {
		if fname[i] == '.' {
			return fname[i:]
		}
	}
	return ""
}

// isText проверяет, является ли содержимое текстовым
func isText(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return false
		}
		if b < 32 && b != 9 && b != 10 && b != 13 {
			return false
		}
	}
	return len(data) > 0
}
