package main

import (
	"os"

	digitalprofileHandlers "visiologyDataUpdate/internal/digital_profile/handlers"
	digitalprofileToken "visiologyDataUpdate/internal/digital_profile/token"
	visiologyHandlers "visiologyDataUpdate/internal/visiology/handlers"
	visiologyToken "visiologyDataUpdate/internal/visiology/token"

	"visiologyDataUpdate/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	DigitalProfileURL    string
	DigitalProfileBearer string
	VisiologyURL         string
	VisiologyBearer      string
	VisiologyAPIVersion  string
}

// loadEnv загружает переменные окружения и инициализирует конфигурацию.
func loadEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	digitalProfileURL := os.Getenv("DIGITAL_PROFILE_BASE_URL")
	digitalProfileBearer, err := digitalprofileToken.GetToken(digitalProfileURL)
	if err != nil {
		return nil, err
	}

	visiologyURL := os.Getenv("VISIOLOGY_BASE_URL")
	visiologyBearer, err := visiologyToken.GetToken(visiologyURL)
	if err != nil {
		return nil, err
	}

	visiologyAPIVersion := os.Getenv("VISIOLOGY_API_VERSION")

	return &Config{
		DigitalProfileURL:    digitalProfileURL,
		DigitalProfileBearer: "Bearer " + digitalProfileBearer,
		VisiologyURL:         visiologyURL,
		VisiologyBearer:      "Bearer " + visiologyBearer,
		VisiologyAPIVersion:  visiologyAPIVersion,
	}, nil
}

// main является точкой входа в программу.
func main() {
	// Инициализируем логгер
	logger.InitLogger()
	config, err := loadEnv()
	if err != nil {
		logger.Fatal("Ошибка загрузки файла .env: ", err)
	}

	// Проверяем, установлена ли переменная DEBUG в "True"
	if os.Getenv("DEBUG") == "True" {
		logger.Info("Отладка включена: Извлечение конфигурации")
		logger.Info("URL цифрового профиля: ", config.DigitalProfileURL)
		logger.Info("Bearer цифрового профиля: ", config.DigitalProfileBearer)
		logger.Info("URL Visiology: ", config.VisiologyURL)
		logger.Info("Bearer Visiology: ", config.VisiologyBearer)
	}

	logger.Info("Получение ответа от API цифрового профиля")

	// Получение ответа от API цифрового профиля
	digitalProfileResponse := digitalprofileHandlers.GetHandler(config.DigitalProfileURL, config.DigitalProfileBearer)

	// Отправка ответа на платформу Visiology
	visiologyHandlers.PostHandler(digitalProfileResponse, config.VisiologyURL, config.VisiologyAPIVersion, config.VisiologyBearer)

	logger.Info("Программа завершена")
}
