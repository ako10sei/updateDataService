package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	digitalprofile "visiologyDataUpdate/internal/digital_profile/handlers"
	visiology "visiologyDataUpdate/internal/visiology/structs"
)

// OrgIDs формуирует срез id организаций валидных для обработки и обновления по данным ЦП.
var OrgIDs = []int{3, 27, 11, 12, 5, 17, 22, 7, 21, 20, 13, 10, 24, 14, 15, 16, 18, 6, 19, 9, 8, 30, 43}

// maxIterations ограничивает число итераций для обработки организаций.
const maxIterations = 22

// PostHandler обрабатывает ответ от цифрового профиля и отправляет его на платформу Visiology.
// Он создает тело запроса, содержащее данные организации, и отправляет его в виде POST-запроса по указанному URL-адресу Visiology.
//
// Параметры:
// - digitalProfileResponse: Ответ от API цифрового профиля, содержащий данные об организациях.
// - visiologyUrl: URL-адрес платформы Visiology.
// - visiologyApiVersion: Версия API Visiology, которая будет использоваться.
// - visiologyBearer: Токен авторизации для проверки подлинности с платформой Visiology.
//
//nolint:funlen
func PostHandler(
	digitalProfileResponse digitalprofile.GetResponse,
	visiologyURL,
	visiologyAPIVersion,
	visiologyBearer string) {

	// TODO: Реализовать передачу параметров: Проектная мощность, Филиалы.
	// TODO: Сделать логирование всех процессов в работе с цифровым профилем и отправки запросов в Visiology.
	// Инициализация переменных
	var column visiology.Column
	var fields = column.GetAllFields()
	var rownum = 0
	var requestBody []map[string]any
	// Создание тела запроса, содержащего данные организаций
	for rownum != maxIterations+1 {
		for _, org := range digitalProfileResponse.Organizations {
			if rownum > maxIterations {
				break
			}
			// Т.к. необходимо отправлять данные только для указанных идентификаторов организаций (см. OrgIds),
			// Требуется добавить проверку на работу с указанными идентификаторами.
			// В случае, если строится дата для оргнизации, которая не входит в OrgIds,
			// То условие не пропустит данную организацию для построения JSON, который далее отправится в Visiology.
			if org.ID == OrgIDs[rownum] {
				for _, field := range fields {
					rowData := map[string]any{
						"rownum": rownum,
						"values": []map[string]any{
							{
								"column": field,
								"value":  org.GetValueByField()[field],
							},
						},
					}
					requestBody = append(requestBody, rowData)
				}
			} else {
				continue
			}
			rownum++
		}
	}
	// Маршалирование тела запроса в JSON-формате для отправки на сервер Visiology
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	// Создание HTTP-запроса с телом запроса
	req, err := http.NewRequest("POST", visiologyURL+"/update", bytes.NewBuffer(jsonBody))
	if err != nil {
		return
	}
	// Добавление заголовков HTTP-запроса
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", fmt.Sprintf("%d", len(jsonBody)))
	req.Header.Add("Authorization", visiologyBearer)
	req.Header.Add("X-Api-Version", visiologyAPIVersion)
	req.Header.Add("Host", "<calculated when request is sent>")
	// Отправка HTTP-запроса и получение ответа
	client := &http.Client{}
	resp, err := client.Do(req) //nolint:bodyclose
	if err != nil {
		log.Fatal("Ошибка при отправке HTTP-запроса:", err)

		return
	}
	// Закрытие тела ответа после завершения работы с ним
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Ошибка закрытия тела ответа:", err)
		}
	}(resp.Body)

	// Проверка статуса HTTP-ответа
	if resp.StatusCode != http.StatusOK {
		// Чтение тела ответа в случае некорректного статуса HTTP
		bodyBytes, err := io.ReadAll(resp.Body)
		// Вывод статуса HTTP и тела ответа
		fmt.Println("Non-ok HTTP status:", resp.StatusCode)
		fmt.Println("GetResponse body:", string(bodyBytes))
		if err != nil {
			log.Panic("Ошибка во время чтения тела ответа:", err)
		}
	}
}
