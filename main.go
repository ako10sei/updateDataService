package main

import (
	"fmt"
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
	digitalProfileResponse := digitalprofile.Handler(digitalProfileUrl, digitalProfileBearer)
	visiologyResponse := visiology.Handler(visiologyUrl, visiologyBearer, visiologyApiVersion)

	for _, d := range digitalProfileResponse.Organizations {
		for _, v := range visiologyResponse.Values {
			fmt.Printf("Names: %s %s\n", d.Name, v[2])
		}
	}

}
