package main

import (
	"os"
	"visiologyDataUpdate/internal/app" // Импортируем новый пакет
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
