.PHONY: all build clean run test

# Компиляция всех утилит
all: build

build:
	@mkdir -p bin
	@for dir in */; do \
		if [ -f "$$dir/main.go" ]; then \
			echo "Building $$dir..."; \
			go build -o "bin/$${dir%/}" "$$dir/main.go"; \
		fi \
	done

# Очистка
clean:
	@rm -rf bin
	@echo "Очищено"

# Запуск конкретной утилиты
run-%:
	@go run $*/main.go $(ARGS)

# Тестирование
test:
	@./test_all.sh

# Установка в систему
install:
	@sudo cp bin/* /usr/local/bin/

# Показать доступные команды
list:
	@ls -d */ | sed 's/\///'