package main

import (
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	digitalprofile "visiologyDataUpdate/internal/digital_profile/handlers"
	digitalprofiletoken "visiologyDataUpdate/internal/digital_profile/token"
	visiology "visiologyDataUpdate/internal/visiology/handlers"
)

var (
	digitalProfileURL    string
	digitalProfileBearer string
	visiologyURL         string
	visiologyBearer      string
	visiologyAPIVersion  string
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
	digitalProfileURL = os.Getenv("DIGITAL_PROFILE_BASE_URL")
	digitalProfileBearer = "Bearer " + digitalprofiletoken.GetToken(digitalProfileURL)

	// Получение URL-адреса для платформы Visiology и токена доступа
	visiologyURL = os.Getenv("VISIOLOGY_BASE_URL")
	visiologyBearer = "Bearer " + os.Getenv("VISIOLOGY_API_TOKEN")

	// Получение версии API Visiology
	visiologyAPIVersion = os.Getenv("VISIOLOGY_API_VERSION")
}

// main является точкой входа в программу. Она инициализирует необходимые переменные,
// получает данные из API цифрового профиля и отправляет данные на платформу Visiology.
func main() {
	// Инициализация логгера
	logger := slog.Default()
	logger.Info("Получение ответа от API цифрового профиля")
	logger.Info("URL: %s, Bearer: %s", digitalProfileURL, digitalProfileBearer)
	// Получение ответа от API цифрового профиля
	digitalProfileResponse := digitalprofile.GetHandler(digitalProfileURL, digitalProfileBearer, logger)

	// Отправка ответа на платформу Visiology с использованием маркера доступа и версии API
	defer visiology.PostHandler(digitalProfileResponse, visiologyURL, visiologyAPIVersion, visiologyBearer)
}
