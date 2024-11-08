package main

import (
	"os"

	"github.com/ako10sei/updateDataService/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
