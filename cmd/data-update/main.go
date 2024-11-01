package main

import (
	"github.com/ako10sei/updateDataService/internal/app"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
