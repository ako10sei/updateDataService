package organization

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"visiologyDataUpdate/digital_profile/structs"
)

type GetResponse struct {
	Count         int                    `json:"count"`
	Next          any                    `json:"next"`
	Previous      any                    `json:"previous"`
	Organizations []structs.Organization `json:"results"`
}

// GetHandler является функцией, которая отправляет GET-запрос на указанный URL с указанным маркером доступа,
// обрабатывает ответ и возвращает структуру GetResponse, содержащую данные организаций.
//
// Параметры:
// - digitalProfileUrl: Строка, представляющая базовый URL API цифрового профиля.
// - digitalProfileBearer: Строка, представляющая токен доступа для API цифрового профиля.
//
// Возвращаемое значение:
// - GetResponse: Структура, содержащая данные организаций, полученные из ответа API цифрового профиля.
func GetHandler(digitalProfileUrl string, digitalProfileBearer string) GetResponse {
	// Создание нового HTTP-запроса
	req, err := http.NewRequest("GET", digitalProfileUrl+"organizations", nil)
	if err != nil {
		log.Fatal("Ошибка: %v", err)
	}

	// Добавление маркера доступа в заголовок запроса
	req.Header.Add("Authorization", digitalProfileBearer)

	// Создание нового HTTP-клиента
	client := &http.Client{}

	// Отправка HTTP-запроса и получение ответа
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Ошибка в ответе.\n[ERROR] -", err)
	}

	// Закрытие тела ответа после завершения работы с ним
	defer resp.Body.Close()

	// Проверка статуса HTTP-ответа
	if resp.StatusCode != http.StatusOK {
		// Чтение тела ответа в случае некорректного статуса HTTP
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Ошибка во время чтения тела ответа:", err)
		}

		// Вывод статуса HTTP и тела ответа
		fmt.Println("Non-ok HTTP status:", resp.StatusCode)
		fmt.Println("GetResponse body:", string(bodyBytes))

		// Возврат пустой структуры в случае некорректного статуса HTTP
		return GetResponse{}
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка во время обработки JSON:", err)
	}

	// Десериализация тела ответа в структуру GetResponse
	var response GetResponse
	json.Unmarshal(body, &response)

	// Возврат структуры GetResponse
	return response
}
