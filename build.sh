#!/bin/bash

set -e

echo "Начало сборки утилит..."
mkdir -p bin

# Перебор всех поддиректорий в commands/
for dir in */; do
    # Убираем слеш в конце
    dir="${dir%/}"
    
    # Проверяем наличие main.go
    if [ -f "$dir/main.go" ]; then
        echo "Компиляция: $dir/main.go -> bin/$dir"
        
        # Сборка для текущей ОС
        go build -o "bin/$dir" "$dir/main.go"
        
        # Сборка для Windows (опционально)
        # GOOS=windows GOARCH=amd64 go build -o "bin/${dir}.exe" "$dir/main.go"
    else
        echo "Пропуск $dir: main.go не найден"
    fi
done

echo "Сборка успешно завершена."
echo "Исполняемые файлы находятся в папке ./bin/"
ls -lh bin/