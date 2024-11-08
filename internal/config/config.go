package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/ako10sei/updateDataService/internal/token"
)

type Config struct {
	Env                  string
	DigitalProfileURL    string
	DigitalProfileBearer string
	VisiologyURL         string
	VisiologyBearer      string
	VisiologyAPIVersion  string
}

func Load(log *slog.Logger) (*Config, error) {
	errChan := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	digitalProfileProvider := token.NewDigitalProfileTokenProvider(os.Getenv("DIGITAL_PROFILE_BASE_URL"))
	visiologyProvider := token.NewVisiologyTokenProvider(os.Getenv("VISIOLOGY_BASE_URL"))

	var digitalProfileBearer, visiologyBearer string

	// Запуск первой горутины для получения токена Digital Profile
	go func() {
		defer wg.Done()
		var err error
		digitalProfileBearer, err = digitalProfileProvider.GetToken(log)
		if err != nil {
			errChan <- err
		}
	}()

	// Запуск второй горутины для получения токена Visiology
	go func() {
		defer wg.Done()
		var err error
		visiologyBearer, err = visiologyProvider.GetToken(log)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	// Проверяем наличие ошибок
	if err, ok := <-errChan; ok {
		return nil, err
	}

	return &Config{
		Env:                  os.Getenv("ENVIRONMENT"),
		DigitalProfileURL:    os.Getenv("DIGITAL_PROFILE_BASE_URL"),
		DigitalProfileBearer: "Bearer " + digitalProfileBearer,
		VisiologyURL:         os.Getenv("VISIOLOGY_BASE_URL"),
		VisiologyBearer:      "Bearer " + visiologyBearer,
		VisiologyAPIVersion:  os.Getenv("VISIOLOGY_API_VERSION"),
	}, nil
}
