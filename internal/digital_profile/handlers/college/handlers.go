package college

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"updateDataService/internal/digital_profile/handlers/college/structs"
)

// GetResponse представляет ответ от API цифрового профиля.
type GetResponse struct {
	Count         int                    `json:"count"`
	Next          any                    `json:"next"`
	Previous      any                    `json:"previous"`
	Organizations []structs.Organization `json:"results"`
}

// GetHandler отправляет GET-запрос на API цифрового профиля и возвращает структуру GetResponse.
func GetHandler(digitalProfileURL, digitalProfileBearer string, log *slog.Logger) (GetResponse, error) {
	log.Info("Отправка GET-запроса на API цифрового профиля", "url", digitalProfileURL+"organizations")

	req, err := createRequest(digitalProfileURL+"organizations", digitalProfileBearer)
	if err != nil {
		log.Error("Ошибка создания HTTP-запроса", "error", err)
		return GetResponse{}, err
	}

	resp, err := sendRequest(req) //nolint:bodyclose
	if err != nil {
		log.Error("Ошибка при отправке HTTP-запроса", "error", err)
		return GetResponse{}, err
	}
	defer closeResponse(resp.Body, log)

	if resp.StatusCode != http.StatusOK {
		return GetResponse{}, handleNonOKResponse(resp, log)
	}

	response, err := parseResponse(resp.Body)
	if err != nil {
		log.Error("Ошибка при разборе ответа", "error", err)
		return GetResponse{}, err
	}

	log.Info("Ответ получен успешно", "count", response.Count)
	return response, nil
}

// createRequest создает новый HTTP-запрос с авторизационным заголовком.
func createRequest(url, bearer string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	return req, nil
}

// sendRequest отправляет HTTP-запрос и возвращает ответ.
func sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

// handleNonOKResponse обрабатывает некорректный статус HTTP-ответа и возвращает ошибку.
func handleNonOKResponse(resp *http.Response, log *slog.Logger) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка при чтении тела ответа с некорректным статусом", "error", err)
		return err
	}
	log.Error("Некорректный статус HTTP", "status", resp.StatusCode, "body", string(bodyBytes))
	return errors.New("получен некорректный статус ответа от сервера")
}

// parseResponse читает тело ответа и декодирует JSON в структуру GetResponse.
func parseResponse(body io.ReadCloser) (GetResponse, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return GetResponse{}, err
	}

	var response GetResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return GetResponse{}, err
	}
	return response, nil
}

// closeResponse закрывает тело ответа и логирует ошибку, если она произошла.
func closeResponse(body io.ReadCloser, log *slog.Logger) {
	if err := body.Close(); err != nil {
		log.Error("Ошибка закрытия тела ответа", "error", err)
	}
}
