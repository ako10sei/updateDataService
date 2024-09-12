package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	grantType = "client_credentials"
	scope     = "digital_profile"
)

var (
	body []byte
)

// init выполняется при инициализации пакета и загружает переменные окружения.
func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Ошибка загрузки файла .env", "error: ", err)
	}

	body = createRequestBody(grantType,
		os.Getenv("DIGITAL_PROFILE_CLIENT_ID"),
		os.Getenv("DIGITAL_PROFILE_CLIENT_SECRET"),
		scope)
}

func createRequestBody(grantType, clientID, clientSecret, scope string) []byte {
	body := fmt.Sprintf(`{
		"grant_type": "%s",
		"client_id": "%s",
		"client_secret": "%s",
		"scope": "%s"
	}`, grantType, clientID, clientSecret, scope)
	return []byte(body)
}

// Token представляет собой структуру для доступа к токену.
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GetToken получает токен доступа из указанного URL.
func GetToken(digitalProfileURL string, log *slog.Logger) (string, error) {
	log.Debug("Отправка запроса на получение токена", "URL: ", digitalProfileURL)

	req, err := http.NewRequest("POST", digitalProfileURL+"oauth2/token/", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("ошибка создания HTTP-запроса: %w", err)
	}

	resp, err := (&http.Client{}).Do(req) //nolint:bodyclose
	if err != nil {
		return "", fmt.Errorf("ошибка при отправке HTTP-запроса: %w", err)
	}
	defer closeResponse(resp.Body, log)

	if resp.StatusCode != http.StatusOK {
		handleNonOKResponse(resp, log)
		return "", fmt.Errorf("неверный статус ответа: %d", resp.StatusCode)
	}

	var token Token
	// Чтение и десериализация тела ответа
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("ошибка десериализации тела ответа: %w", err)
	}

	log.Info("Токен доступа успешно получен", "accessToken: ", token.AccessToken)
	return token.AccessToken, nil
}

// handleNonOKResponse обрабатывает ошибку сервера в случае ошибки
func handleNonOKResponse(resp *http.Response, log *slog.Logger) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка во время чтения тела ответа", "error: ", err)
		return
	}
	log.Error("Некорректный статус HTTP", "status: ", resp.StatusCode, "body: ", string(bodyBytes))
}

// closeResponse закрывает тело ответа и логирует ошибку, если она произошла.
func closeResponse(body io.ReadCloser, log *slog.Logger) {
	if err := body.Close(); err != nil {
		log.Error("Ошибка закрытия тела ответа", "error: ", err)
	}
}
