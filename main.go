package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	digitalprofile "visiologyDataUpdate/digital_profile/rest/organization"
	visiology "visiologyDataUpdate/visiology/rest/organization"
)

var (
	digitalProfileUrl    string
	digitalProfileBearer string
	visiologyUrl         string
	visiologyBearer      string
	visiologyApiVersion  string
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
	visiologyApiVersion = os.Getenv("VISIOLOGY_API_VERSION")

}

func main() {
	digitalProfileResponse := digitalprofile.GetHandler(digitalProfileUrl, digitalProfileBearer)
	visiologyResponse := visiology.GetHandler(visiologyUrl, visiologyBearer, visiologyApiVersion)

	defer visiology.PostHandler(digitalProfileResponse, visiologyResponse, visiologyUrl, visiologyApiVersion, visiologyBearer)

}
