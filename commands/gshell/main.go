.PHONY: all build clean install run help

# Цвета для вывода
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

# Список всех утилит
UTILS = arch cat cd clear cp date df du file free head hexdump history \
        kill ls mkdir nl ps pwd pwgen rm rmdir tail tar touch uname unzip wc zip

# Оболочка
SHELL_UTIL = gshell

all: build

# Сборка всех утилит
build:
	@echo "$(GREEN)Сборка утилит...$(NC)"
	@mkdir -p bin
	@for util in $(UTILS); do \
		echo "$(YELLOW)  Сборка $$util...$(NC)"; \
		(cd $$util && go build -o ../bin/$$util .) || exit 1; \
	done
	@echo "$(YELLOW)  Сборка оболочки...$(NC)"
	@(cd $(SHELL_UTIL) && go build -o ../bin/$(SHELL_UTIL) .) || exit 1
	@chmod +x bin/*
	@echo "$(GREEN)Готово!$(NC)"

# Сборка с отладочной информацией
build-debug:
	@echo "$(GREEN)Сборка с отладкой...$(NC)"
	@mkdir -p bin
	@for util in $(UTILS); do \
		echo "$(YELLOW)  Сборка $$util...$(NC)"; \
		(cd $$util && go build -gcflags="all=-N -l" -o ../bin/$$util .) || exit 1; \
	done
	@(cd $(SHELL_UTIL) && go build -gcflags="all=-N -l" -o ../bin/$(SHELL_UTIL) .) || exit 1
	@chmod +x bin/*
	@echo "$(GREEN)Готово!$(NC)"

# Сборка для Windows
build-windows:
	@echo "$(GREEN)Сборка для Windows...$(NC)"
	@mkdir -p bin/windows
	@for util in $(UTILS); do \
		echo "$(YELLOW)  Сборка $$util.exe...$(NC)"; \
		(cd $$util && GOOS=windows GOARCH=amd64 go build -o ../bin/windows/$$util.exe .) || exit 1; \
	done
	@(cd $(SHELL_UTIL) && GOOS=windows GOARCH=amd64 go build -o ../bin/windows/$(SHELL_UTIL).exe .) || exit 1
	@echo "$(GREEN)Готово!$(NC)"

# Сборка для Linux
build-linux:
	@echo "$(GREEN)Сборка для Linux...$(NC)"
	@mkdir -p bin/linux
	@for util in $(UTILS); do \
		echo "$(YELLOW)  Сборка $$util...$(NC)"; \
		(cd $$util && GOOS=linux GOARCH=amd64 go build -o ../bin/linux/$$util .) || exit 1; \
	done
	@(cd $(SHELL_UTIL) && GOOS=linux GOARCH=amd64 go build -o ../bin/linux/$(SHELL_UTIL) .) || exit 1
	@chmod +x bin/linux/*
	@echo "$(GREEN)Готово!$(NC)"

# Сборка для macOS
build-macos:
	@echo "$(GREEN)Сборка для macOS...$(NC)"
	@mkdir -p bin/macos
	@for util in $(UTILS); do \
		echo "$(YELLOW)  Сборка $$util...$(NC)"; \
		(cd $$util && GOOS=darwin GOARCH=amd64 go build -o ../bin/macos/$$util .) || exit 1; \
	done
	@(cd $(SHELL_UTIL) && GOOS=darwin GOARCH=amd64 go build -o ../bin/macos/$(SHELL_UTIL) .) || exit 1
	@chmod +x bin/macos/*
	@echo "$(GREEN)Готово!$(NC)"

# Установка в систему
install: build
	@echo "$(GREEN)Установка в /usr/local/bin...$(NC)"
	@sudo cp bin/* /usr/local/bin/
	@echo "$(GREEN)Установка завершена!$(NC)"

# Запуск оболочки
run: build
	@echo "$(GREEN)Запуск Go Shell...$(NC)"
	@./bin/$(SHELL_UTIL)

# Запуск с конкретной командой
run-cmd: build
	@echo "$(GREEN)Запуск команды...$(NC)"
	@./bin/$(SHELL_UTIL) -c "$(CMD)"

# Очистка
clean:
	@echo "$(YELLOW)Очистка...$(NC)"
	@rm -rf bin/
	@for util in $(UTILS); do \
		rm -f $$util/$$util; \
	done
	@rm -f $(SHELL_UTIL)/$(SHELL_UTIL)
	@echo "$(GREEN)Очищено!$(NC)"

# Форматирование кода
fmt:
	@echo "$(YELLOW)Форматирование кода...$(NC)"
	@for util in $(UTILS); do \
		(cd $$util && go fmt .) || exit 1; \
	done
	@(cd $(SHELL_UTIL) && go fmt .) || exit 1
	@echo "$(GREEN)Готово!$(NC)"

# Справка
help:
	@echo "$(GREEN)Доступные команды:$(NC)"
	@echo "  make build          - собрать все утилиты"
	@echo "  make build-debug    - собрать с отладкой"
	@echo "  make build-windows  - собрать для Windows"
	@echo "  make build-linux    - собрать для Linux"
	@echo "  make build-macos    - собрать для macOS"
	@echo "  make install        - установить в систему"
	@echo "  make run            - запустить оболочку"
	@echo "  make run-cmd CMD='ls -l' - выполнить команду"
	@echo "  make clean          - очистить сборку"
	@echo "  make fmt            - отформатировать код"
	@echo "  make help           - показать эту справку"