#!/bin/bash

set -e

echo "Начало сборки утилит..."
mkdir -p bin

# Перебор всех поддиректорий в commands/
for dir in commands/*/; do
    # Убираем слеш в конце
    dir="${dir%/}"
   
    # Проверяем наличие main.go
    if [ -f "$dir/main.go" ]; then
        # Извлекаем имя утилиты (последнюю часть пути)
        util_name="${dir##*/}"
        echo "Компиляция: $dir/main.go -> bin/$util_name"
       
        # Сборка для текущей ОС
        go build -o "bin/$util_name" "$dir/main.go"
    else
        echo "Пропуск $dir: main.go не найден"
    fi
done

echo "Сборка успешно завершена."
echo "Исполняемые файлы находятся в папке ./bin/"
ls -lh bin/
