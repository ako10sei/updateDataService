package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"visiologyDataUpdate/digital_profile/rest/organization"
)

var (
	digitalProfileUrl    string
	digitalProfileBearer string
	visiologyUrl         string
	visiologyBearer      string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки.env файла")
	}
	// Получение DIGITAL_PROFILE_BASE_URL и DIGITAL_PROFILE_API_TOKEN
	digitalProfileUrl = os.Getenv("DIGITAL_PROFILE_BASE_URL")
	digitalProfileBearer = "Bearer " + os.Getenv("DIGITAL_PROFILE_API_TOKEN")

	// Получение VISIOLOGY_BASE_URL и VISIOLOGY_API_TOKEN
	visiologyUrl = os.Getenv("VISIOLOGY_BASE_URL")
	visiologyBearer = "Bearer " + os.Getenv("VISIOLOGY_API_TOKEN")

}

func main() {
	digitalProfileResponse := organization.Handler(digitalProfileUrl, digitalProfileBearer)
	for _, p := range digitalProfileResponse.Organizations {
		fmt.Println("Organization", p.ID, p.Name, p.ShortName, p.Parent)
	}
}
