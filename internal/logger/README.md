# Пакет log

Этот пакет log предоставляет удобный и гибкий инструмент для логирования в Go. Он включает в себя цветные стили для уровней логирования и предоставляет простые функции для логирования сообщений на различных уровнях.

## Использование

После импорта вы можете использовать следующие функции для логирования:

- `log.Info(msg string, args ...any)`: логирует информационное сообщение.
- `log.Debug(msg string, args ...any)`: логирует отладочное сообщение.
- `log.Error(msg string, args ...any)`: логирует сообщение об ошибке.
- `log.Fatal(msg string, args ...any)`: логирует фатальное сообщение и завершает программу с кодом 1.

## Пример использования

```go
package main

import (
"visiologyupdategolang/internal/log"
)

func main() {
log.InitLogger() // Инициализируем логгер

log.Info("This is an informational message")
log.Debug("This is a debug message")
log.Error("This is an error message")
log.Fatal("This is a fatal error message")
}
```

## Настройка

По умолчанию пакет log выводит логи в стандартный поток вывода (stdout). Если вам нужно изменить это поведение, вы можете настроить логгер с помощью следующей функции:

```go
func InitLogger(out *os.File) {
log.Logger = log.New(out, "", 0)
}
```

Вы можете передать любой `*os.File` в качестве аргумента для настройки вывода логов.

## Цветные стили

Пакет log использует пакет `github.com/fatih/color` для цветного вывода уровней логирования. Вы можете настроить цвета для уровней логирования, изменив соответствующие переменные:

- `log.InfoColor`: цвет для информационных сообщений.
- `log.DebugColor`: цвет для отладочных сообщений.
- `log.ErrorColor`: цвет для сообщений об ошибках.
- `log.FatalColor`: цвет для фатальных сообщений.

## Формат сообщений

Сообщения логирования форматируются следующим образом:


```[YYYY-MM-DD HH:MM:SS] [LEVEL] MESSAGE: [ARGUMENTS...]```


- `YYYY-MM-DD HH:MM:SS`: текущая временная метка.
- `LEVEL`: уровень логирования (INFO, DEBUG, ERROR, FATAL).
- `MESSAGE`: основное сообщение логирования.
- `[ARGUMENTS...]`: дополнительные аргументы, переданные в функцию логирования.
