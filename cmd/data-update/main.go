package main

import (
	"os"
	conf "visiologyDataUpdate/internal/config"
	digitalprofileHandlers "visiologyDataUpdate/internal/digital_profile/handlers"
	"visiologyDataUpdate/internal/logger"
	visiologyHandlers "visiologyDataUpdate/internal/visiology/handlers"
)

// main является точкой входа в программу.
func main() {
	// Инициализируем логгер
	log := logger.SetupLogger("local")
	cfg, err := conf.Load(log)
	if err != nil {
		log.Error("Ошибка загрузки файла .env: ", "error", err)
	}

	// Проверяем, установлена ли переменная DEBUG в "True"
	if os.Getenv("DEBUG") == "True" {
		log.Info("Отладка включена: Извлечение конфигурации")
		log.Info("URL цифрового профиля: ", "url", cfg.DigitalProfileURL)
		log.Info("Bearer цифрового профиля: ", "token", cfg.DigitalProfileBearer)
		log.Info("URL Visiology: ", "url", cfg.VisiologyURL)
		log.Info("Bearer Visiology: ", "token", cfg.VisiologyBearer)
	}

	log.Info("Получение ответа от API цифрового профиля")

	// Получение ответа от API цифрового профиля
	digitalProfileResponse := digitalprofileHandlers.GetHandler(cfg.DigitalProfileURL, cfg.DigitalProfileBearer, log)

	// Отправка ответа на платформу Visiology
	visiologyHandlers.PostHandler(digitalProfileResponse, cfg.VisiologyURL, cfg.VisiologyAPIVersion, cfg.VisiologyBearer, log)

	log.Info("Программа завершена!")
}
