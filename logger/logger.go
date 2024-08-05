package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

var Logger *log.Logger

// Инициализируем цветные стили для уровней логирования.
var (
	InfoColor  = color.New(color.FgBlue).SprintFunc()
	DebugColor = color.New(color.FgGreen).SprintFunc()
	ErrorColor = color.New(color.FgRed).SprintFunc()
	FatalColor = color.New(color.FgMagenta).SprintFunc()
)

// InitLogger инициализирует логгер с текстовым форматом.
func InitLogger() {
	Logger = log.New(os.Stdout, "", 0)
}

// getCurrentTime возвращает текущую временную метку в нужном формате.
func getCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Info logs a message at the Info level.
func Info(msg string, args ...any) {
	Logger.Println(formatMessage("INFO", InfoColor, msg, args...))
}

// Debug logs a message at the Debug level.
func Debug(msg string, args ...any) {
	Logger.Println(formatMessage("DEBUG", DebugColor, msg, args...))
}

// Error logs a message at the Error level.
func Error(msg string, args ...any) {
	Logger.Println(formatMessage("ERROR", ErrorColor, msg, args...))
}

// Fatal logs a message at the Fatal level and then calls os.Exit(1).
func Fatal(msg string, args ...any) {
	Logger.Println(formatMessage("FATAL", FatalColor, msg, args...))
	os.Exit(1) // Завершаем программу с кодом 1
}

// formatMessage форматирует сообщение для логирования.
func formatMessage(level string, colorFunc func(...any) string, msg string, args ...any) string {
	timeStamp := getCurrentTime()
	if len(args) > 0 {
		return fmt.Sprintf("[%s] [%s] %s: %s", timeStamp, level, colorFunc(msg), fmt.Sprint(args...))
	}
	return fmt.Sprintf("[%s] [%s] %s", timeStamp, level, colorFunc(msg))
}
