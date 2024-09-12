package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const (
	grantType    = "password"
	scope        = "openid profile email roles viqube_api viqubeadmin_api core_logic_facade dashboards_export_service script_service migration_service_api data_collection" //nolint:lll
	responseType = "id_token token"
)

var params = url.Values{
	"grant_type":    {grantType},
	"scope":         {scope},
	"response_type": {responseType},
	"username":      {},
	"password":      {},
}

func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Ошибка загрузки файла .env", "error", err)
	}

	params.Set("username", os.Getenv("VISIOLOGY_USERNAME"))
	params.Set("password", os.Getenv("VISIOLOGY_PASSWORD"))
}

// Token представляет токен доступа.
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GetToken получает токен доступа из указанного URL.
func GetToken(visiologyURL string, log *slog.Logger) (string, error) {
	req, err := http.NewRequest("POST", visiologyURL+"idsrv/connect/token", bytes.NewBufferString(params.Encode()))
	if err != nil {
		return "", fmt.Errorf("ошибка создания HTTP-запроса: %w", err)
	}

	// Заголовок авторизации не утечка! Данный токен един и указан в самой документации к апи (открытый источник)
	// https://visiology-doc.atlassian.net/wiki/spaces/v34/pages/214077490
	req.Header.Add("Authorization", "Basic cHVibGljX3JvX2NsaWVudDpAOVkjbmckXXU+SF4zajY=")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req) //nolint:bodyclose
	if err != nil {
		return "", fmt.Errorf("ошибка при отправке HTTP-запроса: %w", err)
	}
	defer closeResponse(resp.Body, log)

	if resp.StatusCode != http.StatusOK {
		handleNonOKResponse(resp, log)
		return "", fmt.Errorf("неверный статус ответа: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		log.Error("Ошибка десериализации тела ответа", "error", err, "body", string(body), "request", req)
		return "", fmt.Errorf("ошибка десериализации тела ответа: %w", err)
	}

	return token.AccessToken, nil
}

// closeResponse закрывает тело ответа и логирует ошибку, если она произошла.
func closeResponse(body io.ReadCloser, log *slog.Logger) {
	if err := body.Close(); err != nil {
		log.Error("Ошибка закрытия тела ответа", "error", err)
	}
}

// handleNonOKResponse обрабатывает ошибку сервера в случае ошибки
func handleNonOKResponse(resp *http.Response, log *slog.Logger) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка во время чтения тела ответа", "error", err)
		return
	}
	log.Error("Некорректный статус HTTP", "status", resp.StatusCode, "body", string(bodyBytes))
}
