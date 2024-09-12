package config

import (
	"log/slog"
	"os"
	digitalprofileToken "visiologyDataUpdate/internal/digital_profile/token"
	visiologyToken "visiologyDataUpdate/internal/visiology/token"

	"github.com/joho/godotenv"
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
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	env := os.Getenv("ENVIRONMENT")
	digitalProfileURL := os.Getenv("DIGITAL_PROFILE_BASE_URL")
	digitalProfileBearer, err := digitalprofileToken.GetToken(digitalProfileURL, log)
	if err != nil {
		return nil, err
	}

	visiologyURL := os.Getenv("VISIOLOGY_BASE_URL")
	visiologyBearer, err := visiologyToken.GetToken(visiologyURL, log)
	if err != nil {
		return nil, err
	}

	visiologyAPIVersion := os.Getenv("VISIOLOGY_API_VERSION")

	return &Config{
		Env:                  env,
		DigitalProfileURL:    digitalProfileURL,
		DigitalProfileBearer: "Bearer " + digitalProfileBearer,
		VisiologyURL:         visiologyURL,
		VisiologyBearer:      "Bearer " + visiologyBearer,
		VisiologyAPIVersion:  visiologyAPIVersion,
	}, nil
}
