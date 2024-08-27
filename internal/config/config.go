package config

import (
	"os"
	digitalprofileToken "visiologyDataUpdate/internal/digital_profile/token"
	visiologyToken "visiologyDataUpdate/internal/visiology/token"

	"github.com/joho/godotenv"
)

type Config struct {
	DigitalProfileURL    string
	DigitalProfileBearer string
	VisiologyURL         string
	VisiologyBearer      string
	VisiologyAPIVersion  string
}

func Load() (*Config, error) {
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
