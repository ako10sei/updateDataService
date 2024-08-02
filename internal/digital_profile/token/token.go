package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func init() {
	err := godotenv.Load()
	if err != nil {
		// Вывод ошибки и завершение программы, если файл .env не удалось загрузить
		log.Fatal("Ошибка загрузки файла .env")
	}
	body = []byte(fmt.Sprintf(`{
    "grant_type": "%s",
    "client_id": "%s",
    "client_secret": "%s",
    "scope": "%s"
}`,
		grantType,
		os.Getenv("DIGITAL_PROFILE_CLIENT_ID"),
		os.Getenv("DIGITAL_PROFILE_CLIENT_SECRET"),
		scope))
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GetToken получает токен доступа из указанного URL.
// Функция отправляет POST-запрос на указанный URL с необходимыми параметрами,
// включая идентификатор клиента, секрет клиента и область. Затем она читает тело ответа,
// десериализует его в структуру Token и возвращает токен доступа.
//
// Если при отправке HTTP-запроса или чтении тела ответа возникает ошибка,
// функция выводит сообщение об ошибке и завершает работу с ошибкой.
//
// Если HTTP-ответ имеет статус, отличный от 200 (OK), функция читает тело ответа,
// выводит статус HTTP и тело ответа, а затем завершает работу с ошибкой.
func GetToken(digitalProfileURL string) string {
	log.Println("Body до обработки:", string(body))
	req, err := http.NewRequest("POST", digitalProfileURL+"oauth2/token/", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Ошибка получения токена")
	}
	// Создание нового HTTP-клиента
	client := &http.Client{}

	// Отправка HTTP-запроса и получение ответа
	resp, err := client.Do(req) //nolint:bodyclose
	if err != nil {
		log.Fatal("Ошибка при отправке HTTP-запроса:", "error", err)
	}

	// Закрытие тела ответа после завершения работы с ним
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Ошибка закрытия тела ответа:", "error", err)
		}
	}(resp.Body)

	// Проверка статуса HTTP-ответа
	if resp.StatusCode != http.StatusOK {
		// Чтение тела ответа в случае некорректного статуса HTTP
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Panic("Ошибка во время чтения тела ответа:", "error", err)
		}
		// Вывод статуса HTTP и тела ответа
		fmt.Println("Non-ok HTTP status:", resp.StatusCode)
		fmt.Println("GetResponse body:", string(bodyBytes))
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic("Ошибка во время чтения тела ответа:", "error", err)
	}

	var token Token
	// Десериализация тела ответа в структуру
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Println("Request:", req)
		log.Println("GetResponse body:", string(body))
		log.Println("digital_profile_token")
		log.Panic("Ошибка десериализации тела ответа:", "error", err)
	}
	return token.AccessToken
}
