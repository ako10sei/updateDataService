package app

import (
	"log/slog"
	"os"
	"sync"
	"visiologyDataUpdate/external/lib/logger"
	"visiologyDataUpdate/internal/config"
	dp "visiologyDataUpdate/internal/digital_profile/handlers"
	vg "visiologyDataUpdate/internal/visiology/handlers"
)

func Run() error {
	// Инициализация логгера
	log := logger.SetupLogger("local")

	// Загрузка конфигурации
	cfg, err := config.Load(log)
	if err != nil {
		log.Error("Ошибка загрузки файла конфигурации .env", "error", err)
		return err
	}

	// Проверка режима отладки
	if os.Getenv("DEBUG") == "True" {
		log.Info("Отладочный режим включен")
		log.Debug("Конфигурация",
			"DigitalProfileURL", cfg.DigitalProfileURL,
			"DigitalProfileBearer", cfg.DigitalProfileBearer,
			"VisiologyURL", cfg.VisiologyURL,
			"VisiologyBearer", cfg.VisiologyBearer,
		)
	}

	// Создаем экземпляр HandlerConfig для обработки запросов к Visiology
	handlerConfig := vg.NewHandlerConfig(cfg.VisiologyURL, cfg.VisiologyAPIVersion, cfg.VisiologyBearer, log)

	// Получение данных и отправка их на Visiology
	if err := fetchAndSendData(cfg, log, handlerConfig); err != nil {
		return err
	}

	log.Info("Программа завершена успешно!")
	return nil
}

func fetchAndSendData(cfg *config.Config, log *slog.Logger, handlerConfig *vg.HandlerConfig) error {
	// Создаем каналы для передачи данных и ошибок
	digitalProfileCh := make(chan dp.GetResponse) // Замените GetResponse на фактический тип
	errorCh := make(chan error)
	var wg sync.WaitGroup

	// Запускаем goroutine для получения данных из API цифрового профиля
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Получение данных из API цифрового профиля...")
		digitalProfileResponse, err := dp.GetHandler(cfg.DigitalProfileURL, cfg.DigitalProfileBearer, log)
		if err != nil {
			errorCh <- err
			return
		}
		digitalProfileCh <- digitalProfileResponse
	}()

	// Запускаем goroutine для обработки данных после их получения
	go func() {
		wg.Wait()
		close(digitalProfileCh)
	}()

	// Проверка наличия данных
	var digitalProfileResponse dp.GetResponse
	select {
	case digitalProfileResponse = <-digitalProfileCh:
		if len(digitalProfileResponse.Organizations) == 0 {
			log.Error("Получены пустые данные от API цифрового профиля, отправка данных на Visiology не будет выполнена.")
			return nil
		}
	case err := <-errorCh:
		return err
	}

	// Отправка данных на платформу Visiology в отдельной goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Отправка данных на Visiology...")
		if err := handlerConfig.PostHandler(digitalProfileResponse); err != nil {
			errorCh <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errorCh)
	}()

	// Проверяем на ошибки
	if err := <-errorCh; err != nil {
		log.Error("Произошла ошибка при отправке данных на Visiology", "error", err)
		return err
	}

	return nil
}
