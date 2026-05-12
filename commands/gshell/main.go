package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Shell структура оболочки
type Shell struct {
	HomeDir     string
	CurrentDir  string
	PrevDir     string
	History     []string
	CommandsDir string
}

func main() {
	// Флаги оболочки
	var (
		help    = flag.Bool("h", false, "показать справку")
		command = flag.String("c", "", "выполнить команду и выйти")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Создаём оболочку
	shell := &Shell{}
	if err := shell.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации: %v\n", err)
		os.Exit(1)
	}

	// Если указана команда через -c
	if *command != "" {
		if err := shell.ExecuteCommand(*command); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Запускаем интерактивный режим
	shell.Run()
}

// Init инициализирует оболочку
func (s *Shell) Init() error {
	var err error

	// Домашняя директория
	s.HomeDir, err = os.UserHomeDir()
	if err != nil {
		s.HomeDir = "."
	}

	// Текущая директория
	s.CurrentDir, err = os.Getwd()
	if err != nil {
		s.CurrentDir = "."
	}
	s.PrevDir = s.CurrentDir

	// Директория с командами
	s.CommandsDir, err = filepath.Abs(filepath.Join("..", "commands"))
	if err != nil {
		s.CommandsDir = "./commands"
	}

	// Проверяем, существует ли директория команд
	if _, err := os.Stat(s.CommandsDir); os.IsNotExist(err) {
		s.CommandsDir = "../commands"
		if _, err := os.Stat(s.CommandsDir); err != nil {
			fmt.Fprintf(os.Stderr, "Предупреждение: директория команд не найдена\n")
		}
	}

	// Загружаем историю
	s.loadHistory()

	return nil
}

// Run запускает интерактивную оболочку
func (s *Shell) Run() {
	fmt.Println("Go Shell v1.0")
	fmt.Println("Введите 'help' для справки, 'exit' для выхода")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Вывод приглашения
		s.printPrompt()

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Сохраняем в историю
		s.addToHistory(input)

		// Обработка спецсимволов
		expanded := s.expandHistory(input)
		if expanded == "" {
			continue
		}

		// Выполняем команду
		if err := s.ExecuteCommand(expanded); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		}
	}
}

// printPrompt выводит приглашение
func (s *Shell) printPrompt() {
	dir := s.CurrentDir

	// Заменяем домашнюю директорию на ~
	if strings.HasPrefix(dir, s.HomeDir) {
		dir = "~" + strings.TrimPrefix(dir, s.HomeDir)
	}
	if dir == "" {
		dir = "/"
	}

	fmt.Printf("%s $ ", dir)
}

// expandHistory раскрывает !! и !n
func (s *Shell) expandHistory(input string) string {
	// Обработка !!
	if input == "!!" {
		if len(s.History) == 0 {
			fmt.Fprintln(os.Stderr, "!!: нет команд в истории")
			return ""
		}
		lastCmd := s.History[len(s.History)-1]
		fmt.Printf("Выполняется: %s\n", lastCmd)
		return lastCmd
	}

	// Обработка !n
	if strings.HasPrefix(input, "!") && len(input) > 1 {
		numStr := input[1:]
		var num int
		_, err := fmt.Sscanf(numStr, "%d", &num)
		if err != nil {
			return input
		}

		if num < 1 || num > len(s.History) {
			fmt.Fprintf(os.Stderr, "!%d: команда не найдена\n", num)
			return ""
		}

		cmd := s.History[num-1]
		fmt.Printf("Выполняется: %s\n", cmd)
		return cmd
	}

	return input
}

// addToHistory добавляет команду в историю
func (s *Shell) addToHistory(cmd string) {
	if cmd == "" {
		return
	}

	// Не добавляем дубликаты подряд
	if len(s.History) > 0 && s.History[len(s.History)-1] == cmd {
		return
	}

	s.History = append(s.History, cmd)

	// Ограничиваем историю 1000 командами
	if len(s.History) > 1000 {
		s.History = s.History[1:]
	}

	s.saveHistory()
}

// ExecuteCommand выполняет команду
func (s *Shell) ExecuteCommand(input string) error {
	// Разбираем команду и аргументы
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	cmdName := parts[0]
	args := parts[1:]

	// Встроенные команды оболочки
	switch cmdName {
	case "exit", "quit":
		fmt.Println("До свидания!")
		os.Exit(0)

	case "help":
		s.showHelp()
		return nil

	case "history":
		s.showHistory()
		return nil

	case "pwd":
		fmt.Println(s.CurrentDir)
		return nil

	case "cd":
		return s.cdCommand(args)

	case "export":
		return s.exportCommand(args)
	}

	// Внешние команды из папки commands
	return s.runExternalCommand(cmdName, args)
}

// runExternalCommand запускает внешнюю команду из папки commands
func (s *Shell) runExternalCommand(cmdName string, args []string) error {
	// Поиск исполняемого файла
	cmdPath := s.findCommand(cmdName)
	if cmdPath == "" {
		return fmt.Errorf("команда не найдена: %s", cmdName)
	}

	// Создаём команду
	cmd := exec.Command(cmdPath, args...)
	cmd.Dir = s.CurrentDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Устанавливаем переменные окружения
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PWD=%s", s.CurrentDir),
		fmt.Sprintf("OLDPWD=%s", s.PrevDir),
	)

	// Выполняем команду
	err := cmd.Run()

	// Для команды cd нужно обновить состояние оболочки
	if cmdName == "cd" {
		s.updateCurrentDir()
	}

	return err
}

// findCommand ищет команду в директории commands
func (s *Shell) findCommand(cmdName string) string {
	// Пути для поиска
	searchPaths := []string{
		filepath.Join(s.CommandsDir, cmdName, cmdName),
		filepath.Join(s.CommandsDir, cmdName, "main"),
		filepath.Join(s.CommandsDir, cmdName),
		filepath.Join(s.CurrentDir, cmdName),
	}

	// Для Windows добавляем .exe
	if runtime.GOOS == "windows" {
		extensions := []string{".exe", ".bat", ".cmd"}
		newPaths := make([]string, 0)
		for _, p := range searchPaths {
			newPaths = append(newPaths, p)
			for _, ext := range extensions {
				newPaths = append(newPaths, p+ext)
			}
		}
		searchPaths = newPaths
	}

	// Поиск исполняемого файла
	for _, path := range searchPaths {
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			if runtime.GOOS != "windows" && info.Mode()&0111 == 0 {
				continue
			}
			return path
		}
	}

	return ""
}

// cdCommand обрабатывает команду cd
func (s *Shell) cdCommand(args []string) error {
	target := ""

	if len(args) == 0 {
		target = s.HomeDir
	} else if args[0] == "-" {
		target = s.PrevDir
	} else {
		target = args[0]
	}

	// Проверяем существование
	info, err := os.Stat(target)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s: не является директорией", target)
	}

	// Сохраняем предыдущую директорию
	s.PrevDir = s.CurrentDir

	// Меняем директорию
	if err := os.Chdir(target); err != nil {
		return err
	}

	s.CurrentDir = target
	return nil
}

// updateCurrentDir обновляет текущую директорию
func (s *Shell) updateCurrentDir() {
	newDir, err := os.Getwd()
	if err == nil && newDir != s.CurrentDir {
		s.PrevDir = s.CurrentDir
		s.CurrentDir = newDir
	}
}

// exportCommand устанавливает переменную окружения
func (s *Shell) exportCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("export: требуется переменная")
	}

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("export: неверный формат '%s'", arg)
		}

		name := parts[0]
		value := parts[1]

		if err := os.Setenv(name, value); err != nil {
			return err
		}
	}

	return nil
}

// showHistory показывает историю команд
func (s *Shell) showHistory() {
	if len(s.History) == 0 {
		fmt.Println("История пуста")
		return
	}

	fmt.Println("")
	fmt.Println("  N   Команда")
	fmt.Println("--- ----------------------------------------")

	for i, cmd := range s.History {
		fmt.Printf(" %2d  %s\n", i+1, cmd)
	}
	fmt.Println("")
}

// showHelp показывает справку
func (s *Shell) showHelp() {
	fmt.Println("")
	fmt.Println("Go Shell - Доступные команды")
	fmt.Println("")
	fmt.Println("Встроенные команды оболочки:")
	fmt.Println("  cd [директория]    - сменить директорию")
	fmt.Println("  cd -               - вернуться в предыдущую директорию")
	fmt.Println("  pwd                - показать текущую директорию")
	fmt.Println("  history            - показать историю команд")
	fmt.Println("  !!                 - выполнить последнюю команду")
	fmt.Println("  !n                 - выполнить команду под номером n")
	fmt.Println("  export VAR=value   - установить переменную окружения")
	fmt.Println("  help               - показать эту справку")
	fmt.Println("  exit/quit          - выйти из оболочки")
	fmt.Println("")
	fmt.Println("Внешние утилиты (из папки commands/):")
	fmt.Println("  arch, cat, clear, cp, date, df, du, file, find, free")
	fmt.Println("  head, hexdump, kill, ls, mkdir, nl, ps, pwgen, rm, rmdir")
	fmt.Println("  tail, tar, touch, uname, unzip, wc, zip")
	fmt.Println("")
	fmt.Println("Примеры использования:")
	fmt.Println("  $ ls -l -a")
	fmt.Println("  $ cd Documents")
	fmt.Println("  $ cat -n file.txt")
	fmt.Println("  $ mkdir -p new/dir/path")
	fmt.Println("  $ head -n 20 file.txt")
	fmt.Println("")
}

// loadHistory загружает историю из файла
func (s *Shell) loadHistory() {
	historyFile := filepath.Join(s.HomeDir, ".gshell_history")
	file, err := os.Open(historyFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			s.History = append(s.History, line)
		}
	}
}

// saveHistory сохраняет историю в файл
func (s *Shell) saveHistory() {
	historyFile := filepath.Join(s.HomeDir, ".gshell_history")
	file, err := os.Create(historyFile)
	if err != nil {
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, cmd := range s.History {
		fmt.Fprintln(writer, cmd)
	}
	writer.Flush()
}

// printHelp выводит справку по запуску
func printHelp() {
	fmt.Println("")
	fmt.Println("Go Shell (gshell) - Оболочка для запуска утилит")
	fmt.Println("")
	fmt.Println("Использование:")
	fmt.Println("  gshell              - запустить интерактивную оболочку")
	fmt.Println("  gshell -h           - показать эту справку")
	fmt.Println("  gshell -c \"команда\" - выполнить команду и выйти")
	fmt.Println("")
	fmt.Println("Примеры:")
	fmt.Println("  gshell -c \"ls -l -a\"")
	fmt.Println("  gshell -c \"pwd\"")
	fmt.Println("")
}
