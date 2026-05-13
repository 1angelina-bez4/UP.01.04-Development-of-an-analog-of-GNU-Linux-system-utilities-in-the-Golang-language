package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
	"strings"
	"os"
)
const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	special   = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

type PwOptions struct {
	Length      int
	Count       int
	NoSymbols   bool
	NoDigits    bool
	NoUppercase bool
	Secure      bool
}

func main() {
	// Определение флагов
	var (
		help        = flag.Bool("h", false, "показать справку")
		length      = flag.Int("n", 12, "длина пароля")
		count       = flag.Int("c", 1, "количество паролей")
		noSymbols   = flag.Bool("s", false, "без специальных символов")
		noDigits    = flag.Bool("d", false, "без цифр")
		noUppercase = flag.Bool("u", false, "без заглавных букв")
		secure      = flag.Bool("secure", false, "безопасный режим (только буквы и цифры)")
	)
	flag.Parse()

	// Обработка справки
	if *help {
		fmt.Println("pwgen - генерация паролей")
		fmt.Println("Использование: pwgen [-n длина] [-c кол-во] [-s] [-d] [-u] [-secure]")
		fmt.Println("  -n        длина пароля (по умолчанию 12)")
		fmt.Println("  -c        количество паролей (по умолчанию 1)")
		fmt.Println("  -s        без специальных символов")
		fmt.Println("  -d        без цифр")
		fmt.Println("  -u        без заглавных букв")
		fmt.Println("  -secure   безопасный режим (только буквы и цифры)")
		return
	}

	// Проверка длины
	if *length < 4 {
		fmt.Fprintln(os.Stderr, "pwgen: длина пароля должна быть не менее 4")
		os.Exit(1)
	}
	if *length > 128 {
		fmt.Fprintln(os.Stderr, "pwgen: длина пароля не должна превышать 128")
		os.Exit(1)
	}
	if *count < 1 {
		fmt.Fprintln(os.Stderr, "pwgen: количество паролей должно быть не менее 1")
		os.Exit(1)
	}
	if *count > 100 {
		fmt.Fprintln(os.Stderr, "pwgen: количество паролей не должно превышать 100")
		os.Exit(1)
	}

	opts := PwOptions{
		Length:      *length,
		Count:       *count,
		NoSymbols:   *noSymbols,
		NoDigits:    *noDigits,
		NoUppercase: *noUppercase,
		Secure:      *secure,
	}

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Генерируем пароли
	for i := 0; i < opts.Count; i++ {
		password := generatePassword(opts)
		fmt.Println(password)
	}
}

// generatePassword генерирует один пароль
func generatePassword(opts PwOptions) string {
	// Составляем набор символов
	charset := lowercase
	
	if !opts.NoUppercase {
		charset += uppercase
	}
	if !opts.NoDigits && !opts.Secure {
		charset += digits
	}
	if !opts.NoSymbols && !opts.Secure {
		charset += special
	}
	
	if len(charset) == 0 {
		charset = lowercase + digits
	}

	// Генерируем пароль
	password := make([]byte, opts.Length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	// Обеспечиваем наличие хотя бы одного символа из каждой выбранной категории
	if !opts.Secure {
		password = ensureVariety(password, opts)
	}

	return string(password)
}

// ensureVariety обеспечивает разнообразие символов
func ensureVariety(password []byte, opts PwOptions) []byte {
	hasUpper := false
	hasDigit := false
	hasSpecial := false
	
	for _, c := range password {
		if c >= 'A' && c <= 'Z' {
			hasUpper = true
		} else if c >= '0' && c <= '9' {
			hasDigit = true
		} else if strings.ContainsRune(special, rune(c)) {
			hasSpecial = true
		}
	}
	
	// Если чего-то не хватает, добавляем
	if !hasUpper && !opts.NoUppercase {
		password[0] = uppercase[rand.Intn(len(uppercase))]
	}
	if !hasDigit && !opts.NoDigits {
		password[1] = digits[rand.Intn(len(digits))]
	}
	if !hasSpecial && !opts.NoSymbols {
		password[2] = special[rand.Intn(len(special))]
	}
	
	return password
}
