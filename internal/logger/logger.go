package logger

import (
	"log/slog"
	"os"
	"visiologyDataUpdate/internal/config/constants"
	"visiologyDataUpdate/internal/lib/logger/handlers/slogpretty"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case constants.EnvLocal:
		log = setupPrettySlog()
	case constants.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case constants.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
