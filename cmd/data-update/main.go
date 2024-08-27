package main

import (
	"os"
	conf "visiologyDataUpdate/internal/config"
	digitalprofileHandlers "visiologyDataUpdate/internal/digital_profile/handlers"
	"visiologyDataUpdate/internal/log"
	visiologyHandlers "visiologyDataUpdate/internal/visiology/handlers"
)

// main является точкой входа в программу.
func main() {
	// Инициализируем логгер
	log.InitLogger()
	config, err := conf.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env: ", err)
	}

	// Проверяем, установлена ли переменная DEBUG в "True"
	if os.Getenv("DEBUG") == "True" {
		log.Info("Отладка включена: Извлечение конфигурации")
		log.Info("URL цифрового профиля: ", config.DigitalProfileURL)
		log.Info("Bearer цифрового профиля: ", config.DigitalProfileBearer)
		log.Info("URL Visiology: ", config.VisiologyURL)
		log.Info("Bearer Visiology: ", config.VisiologyBearer)
	}

	log.Info("Получение ответа от API цифрового профиля")

	// Получение ответа от API цифрового профиля
	digitalProfileResponse := digitalprofileHandlers.GetHandler(config.DigitalProfileURL, config.DigitalProfileBearer)

	// Отправка ответа на платформу Visiology
	visiologyHandlers.PostHandler(digitalProfileResponse, config.VisiologyURL, config.VisiologyAPIVersion, config.VisiologyBearer)

	log.Info("Программа завершена")
}
