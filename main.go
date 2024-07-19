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

// init является специальной функцией, которая выполняется до функции main.
// Она используется для инициализации переменных, выполнения начальных настроек или выполнения любых других необходимых инициализаций.
func init() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		// Вывод ошибки и завершение программы, если файл .env не удалось загрузить
		log.Fatal("Ошибка загрузки файла .env")
	}

	// Получение URL-адреса для API цифрового профиля и токена доступа
	digitalProfileUrl = os.Getenv("DIGITAL_PROFILE_BASE_URL")
	digitalProfileBearer = "Bearer " + os.Getenv("DIGITAL_PROFILE_API_TOKEN")

	// Получение URL-адреса для платформы Visiology и токена доступа
	visiologyUrl = os.Getenv("VISIOLOGY_BASE_URL")
	visiologyBearer = "Bearer " + os.Getenv("VISIOLOGY_API_TOKEN")

	// Получение версии API Visiology
	visiologyApiVersion = os.Getenv("VISIOLOGY_API_VERSION")
}

// main является точкой входа в программу. Она инициализирует необходимые переменные,
// получает данные из API цифрового профиля и отправляет данные на платформу Visiology.
func main() {
	// Получение ответа от API цифрового профиля
	digitalProfileResponse := digitalprofile.GetHandler(digitalProfileUrl, digitalProfileBearer)

	// Отправка ответа на платформу Visiology с использованием маркера доступа и версии API
	defer visiology.PostHandler(digitalProfileResponse, visiologyUrl, visiologyApiVersion, visiologyBearer)
}