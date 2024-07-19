package organization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	digitalprofile "visiologyDataUpdate/digital_profile/rest/organization"
	"visiologyDataUpdate/visiology/structs"
)

// PostHandler обрабатывает ответ от цифрового профиля и отправляет его на платформу Visiology.
// Он создает тело запроса, содержащее данные организации, и отправляет его в виде POST-запроса по указанному URL-адресу Visiology.
//
// Параметры:
// - digitalProfileResponse: Ответ от API цифрового профиля, содержащий данные об организациях.
// - visiologyUrl: URL-адрес платформы Visiology.
// - visiologyApiVersion: Версия API Visiology, которая будет использоваться.
// - visiologyBearer: Токен авторизации для проверки подлинности с платформой Visiology.
func PostHandler(
	digitalProfileResponse digitalprofile.GetResponse,
	visiologyUrl string,
	visiologyApiVersion string,
	visiologyBearer string) {

	// TODO: Реализовать передачу параметров: Количество студентов общее,
	// TODO: Количество мастеров обучения, Проектная мощность, Филиалы.

	var column structs.Column
	fields := column.GetAllFields()
	// Создание тела запроса, содержащего данные организаций
	requestBody := []map[string]interface{}{}
	for i, org := range digitalProfileResponse.Organizations {
		for _, field := range fields {
			rowData := map[string]interface{}{
				"rownum": i,
				"values": []map[string]interface{}{
					{
						"column": field,
						"value":  org.GetColumnByField()[field],
					},
				},
			}
			requestBody = append(requestBody, rowData)
		}
	}
	// Маршалирование тела запроса в JSON-формате для отправки на сервер Visiology
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}
	// Создание HTTP-запроса с телом запроса
	req, err := http.NewRequest("POST", visiologyUrl+"/update", bytes.NewBuffer(jsonBody))
	if err != nil {
		return
	}
	// Добавление заголовков HTTP-запроса
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(jsonBody)))
	req.Header.Add("Authorization", visiologyBearer)
	req.Header.Add("x-api-version", visiologyApiVersion)
	req.Header.Add("Host", "<calculated when request is sent>")
	// Отправка HTTP-запроса и получение ответа
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Закрытие тела ответа после завершения работы с ним
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Ошибка закрытия тела ответа:", err)
		}
	}(resp.Body)

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
	}
}