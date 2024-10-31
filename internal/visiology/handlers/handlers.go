package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	digitalprofile "visiologyDataUpdate/internal/digital_profile/handlers/college"
	visiology "visiologyDataUpdate/internal/visiology/handlers/college/structs"
)

// HandlerConfig содержит конфигурацию и зависимости для обработки запросов.
type HandlerConfig struct {
	VisiologyURL        string
	VisiologyAPIVersion string
	VisiologyBearer     string
	OrgIDs              []int
	MaxIterations       int
	Log                 *slog.Logger
	Fields              []string
}

// NewHandlerConfig создает новый экземпляр HandlerConfig с заданными параметрами.
func NewHandlerConfig(visiologyURL, visiologyAPIVersion, visiologyBearer string, log *slog.Logger) *HandlerConfig {
	column := visiology.Column{}
	return &HandlerConfig{
		VisiologyURL:        visiologyURL,
		VisiologyAPIVersion: visiologyAPIVersion,
		VisiologyBearer:     visiologyBearer,
		OrgIDs: []int{
			3, 11, 12, 5, 17, 22, 7, 21, 20,
			13, 10, 24, 14, 15, 16, 18, 6, 19,
			9, 8, 30, 43,
		},
		MaxIterations: 23,
		Log:           log,
		Fields:        column.GetAllFields(),
	}
}

// PostHandler обрабатывает ответ от цифрового профиля и отправляет его на платформу Visiology.
func (cfg *HandlerConfig) PostHandler(digitalProfileResponse digitalprofile.GetResponse) error {
	visiologyRequestBody := cfg.createRequestBody(digitalProfileResponse)
	if visiologyRequestBody == nil {
		cfg.Log.Error("Не удалось создать тело запроса, нет валидных организаций для отправки.")
	}

	jsonBody, err := json.MarshalIndent(visiologyRequestBody, "", " ")
	if err != nil {
		cfg.Log.Error("Ошибка при маршалировании JSON тела запроса", "error: ", err)
		return err
	}

	// Проверяем переменную окружения DEBUG
	if os.Getenv("DEBUG") == "True" {
		cfg.Log.Debug("Сформированное тело запроса для Visiology (тестовый режим): ", "jsonBody", string(jsonBody))
		return err
	}

	// Запрашиваем подтверждение перед отправкой
	if !confirmAction("Приступить к обновлению данных Visiology? (Да/Нет): ") {
		cfg.Log.Info("Обновление данных Visiology отменено пользователем.")
		return err
	}

	// Отправляем запрос
	response, err := cfg.sendRequest(jsonBody) //nolint:bodyclose
	if err != nil {
		cfg.Log.Error("Ошибка при отправке HTTP-запроса", "error: ", err)
		return err
	}
	defer closeResponse(response.Body, cfg.Log)

	if response.StatusCode != http.StatusOK {
		handleNonOkResponse(response, cfg.Log)
	} else {
		cfg.Log.Info("Данные успешно отправлены на Visiology")
	}

	return nil
}

// createRequestBody формирует тело запроса на основе ответа от API цифрового профиля.
func (cfg *HandlerConfig) createRequestBody(response digitalprofile.GetResponse) []map[string]any {
	var requestBody []map[string]any

	for rownum, orgID := range cfg.OrgIDs {
		if rownum >= cfg.MaxIterations {
			break
		}

		for _, org := range response.Organizations {
			if org.ID == orgID {
				for _, field := range cfg.Fields {
					rowData := map[string]any{
						"rownum": rownum,
						"values": []map[string]any{{
							"column": field,
							"value":  org.GetValueByField()[field],
						}},
					}
					requestBody = append(requestBody, rowData)
				}
				break
			}
		}
	}

	if len(requestBody) == 0 {
		return nil
	}
	return requestBody
}

// sendRequest отправляет HTTP POST-запрос на указанный URL и возвращает ответ.
func (cfg *HandlerConfig) sendRequest(jsonBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", cfg.VisiologyURL+"viqube/databases/DB/tables/KHV_SPO/records/update", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", cfg.VisiologyBearer)
	req.Header.Add("X-Api-Version", cfg.VisiologyAPIVersion)

	client := &http.Client{}
	return client.Do(req)
}

// closeResponse закрывает тело ответа и логирует ошибку, если она произошла.
func closeResponse(body io.ReadCloser, log *slog.Logger) {
	if err := body.Close(); err != nil {
		log.Error("Ошибка закрытия тела ответа", "error: ", err)
	}
}

// handleNonOkResponse обрабатывает некорректный статус HTTP-ответа.
func handleNonOkResponse(resp *http.Response, log *slog.Logger) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка при чтении тела ответа", "error: ", err)
	}
	log.Error("Некорректный статус HTTP", "status: ", resp.StatusCode, "body: ", string(bodyBytes))
}

// confirmAction запрашивает подтверждение действия у пользователя.
func confirmAction(prompt string) bool {
	fmt.Print(prompt)
	var userResponse string
	_, _ = fmt.Scanln(&userResponse) //nolint:errcheck
	return strings.ToLower(userResponse) == "да"
}
