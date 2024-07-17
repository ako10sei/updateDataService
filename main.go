package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"visiologyDataUpdate/digital_profile"
)

type Response struct {
	Count         int                            `json:"count"`
	Next          any                            `json:"next"`
	Previous      any                            `json:"previous"`
	Organizations []digital_profile.Organization `json:"results"`
}

func main() {
	// Загрузка файла с переменными окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки.env файла")
	}
	// Получение DIGITAL_PROFILE_BASE_URL
	digitalProfileUrl := os.Getenv("DIGITAL_PROFILE_BASE_URL")
	// Получение DIGITAL_PROFILE_API_TOKEN
	var digitalProfileBearer = "Bearer " + os.Getenv("DIGITAL_PROFILE_API_TOKEN")

	// TODO: Все что ниже перенести в соответствуюшие файлы. Файл main будет отвечать за инициализацию и поочередный билд
	// функций обновления.

	req, err := http.NewRequest("GET", digitalProfileUrl, nil)
	req.Header.Add("Authorization", digitalProfileBearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка в ответе.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка во время обработки JSON:", err)
	}

	var response Response

	look := json.Unmarshal(body, &response)

	if look != nil {
		fmt.Println(look)
	}

	for i, p := range response.Organizations {
		fmt.Println("Organization", i+1, ":", p.Name, p.ShortName)
	}
}
