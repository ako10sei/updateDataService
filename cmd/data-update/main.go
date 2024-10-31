package main

import (
	"os"
	"visiologyDataUpdate/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
