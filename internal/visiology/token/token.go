package token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	grantType    = "password"
	scope        = "openid profile email roles viqube_api viqubeadmin_api core_logic_facade dashboards_export_service script_service migration_service_api data_collection"
	responseType = "id_token token"
)

var (
	param = url.Values{}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		// Вывод ошибки и завершение программы, если файл .env не удалось загрузить
		log.Fatal("Ошибка загрузки файла .env")
	}
	param.Set("grant_type", grantType)
	param.Set("scope", scope)
	param.Set("response_type", responseType)
	param.Set("username", os.Getenv("VISIOLOGY_USERNAME"))
	param.Set("password", os.Getenv("VISIOLOGY_PASSWORD"))

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
func GetToken(visiologyURL string) string {

	req, err := http.NewRequest("POST", visiologyURL+"idsrv/connect/token", bytes.NewBufferString(param.Encode()))
	if err != nil {
		log.Fatal("Ошибка получения токена")
	}
	req.Header.Add("Authorization", "Basic cHVibGljX3JvX2NsaWVudDpAOVkjbmckXXU+SF4zajY=")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
			log.Fatal("Ошибка во время чтения тела ответа:", "error", err)
		}
		// Вывод статуса HTTP и тела ответа
		fmt.Println("Non-ok HTTP status:", resp.StatusCode)
		fmt.Println("GetResponse body:", string(bodyBytes))
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка во время чтения тела ответа:", "error", err)
	}

	var token Token
	// Десериализация тела ответа в структуру
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Println("Request:", req)
		log.Println("GetResponse body:", string(body))
		log.Println("visiology_token")
		log.Fatal("Ошибка десериализации тела ответа:", "error", err)
	}
	return token.AccessToken
}
