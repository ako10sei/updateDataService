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

// GetTokenProvider определяет интерфейс для получения токена.
type GetTokenProvider interface {
	GetToken(log *slog.Logger) (string, error)
}

// Token структура для хранения токена доступа.
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// init загружает переменные окружения.
func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Ошибка загрузки файла .env", "error", err)
	}
}

// DigitalProfileTokenProvider реализует TokenProvider для digital_profile.
type DigitalProfileTokenProvider struct {
	URL          string
	ClientID     string
	ClientSecret string
	Scope        string
}

func NewDigitalProfileTokenProvider(url string) *DigitalProfileTokenProvider {
	return &DigitalProfileTokenProvider{
		URL:          url + "oauth2/token/",
		ClientID:     os.Getenv("DIGITAL_PROFILE_CLIENT_ID"),
		ClientSecret: os.Getenv("DIGITAL_PROFILE_CLIENT_SECRET"),
		Scope:        "digital_profile",
	}
}

func (p *DigitalProfileTokenProvider) GetToken(log *slog.Logger) (string, error) {
	log.Debug("Отправка запроса на получение токена Digital Profile", "URL", p.URL)
	body := fmt.Sprintf(`{
		"grant_type": "client_credentials",
		"client_id": "%s",
		"client_secret": "%s",
		"scope": "%s"
	}`, p.ClientID, p.ClientSecret, p.Scope)

	return requestToken(p.URL, body, log, nil)
}

// VisiologyTokenProvider реализует TokenProvider для visiology.
type VisiologyTokenProvider struct {
	URL      string
	Username string
	Password string
	Scope    string
}

func NewVisiologyTokenProvider(url string) *VisiologyTokenProvider {
	return &VisiologyTokenProvider{
		URL:      url + "idsrv/connect/token",
		Username: os.Getenv("VISIOLOGY_USERNAME"),
		Password: os.Getenv("VISIOLOGY_PASSWORD"),
		Scope:    "openid profile email roles viqube_api viqubeadmin_api core_logic_facade dashboards_export_service script_service migration_service_api data_collection",
	}
}

func (p *VisiologyTokenProvider) GetToken(log *slog.Logger) (string, error) {
	log.Debug("Отправка запроса на получение токена Visiology", "URL", p.URL)
	form := url.Values{
		"grant_type":    {"password"},
		"scope":         {p.Scope},
		"response_type": {"id_token token"},
		"username":      {p.Username},
		"password":      {p.Password},
	}

	headers := map[string]string{
		// Заголовок с токеном взят из открытого источника. Не является конфиденциальным ключом доступа.
		"Authorization": "Basic cHVibGljX3JvX2NsaWVudDpAOVkjbmckXXU+SF4zajY=",
		"Content-Type":  "application/x-www-form-urlencoded",
	}

	return requestToken(p.URL, form.Encode(), log, headers)
}

// requestToken выполняет HTTP-запрос для получения токена.
func requestToken(url, body string, log *slog.Logger, headers map[string]string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return "", fmt.Errorf("ошибка создания HTTP-запроса: %w", err)
	}

	// Устанавливаем дополнительные заголовки, если они переданы.
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

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

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", fmt.Errorf("ошибка десериализации тела ответа: %w", err)
	}

	log.Info("Токен успешно получен", "accessToken", token.AccessToken)
	return token.AccessToken, nil
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

// closeResponse закрывает тело ответа и логирует ошибку, если она произошла.
func closeResponse(body io.ReadCloser, log *slog.Logger) {
	if err := body.Close(); err != nil {
		log.Error("Ошибка закрытия тела ответа", "error", err)
	}
}
