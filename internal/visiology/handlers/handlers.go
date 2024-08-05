package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	digitalprofile "visiologyDataUpdate/internal/digital_profile/handlers"
	"visiologyDataUpdate/internal/log"
	visiology "visiologyDataUpdate/internal/visiology/structs"
)

// OrgIDs представляет собой список идентификаторов организаций, валидных для обработки.
var (
	OrgIDs = []int{
		3, 27, 11, 12, 5, 17, 22, 7, 21, 20,
		13, 10, 24, 14, 15, 16, 18, 6, 19,
		9, 8, 30, 43,
	}
	requestBody []map[string]any
	column      visiology.Column
	fields      = column.GetAllFields()
)

// maxIterations ограничивает количество итераций для обработки организаций.
const maxIterations = 23

// PostHandler обрабатывает ответ от цифрового профиля и отправляет его на платформу Visiology.
func PostHandler(
	digitalProfileResponse digitalprofile.GetResponse,
	visiologyURL,
	visiologyAPIVersion,
	visiologyBearer string,
) {
	visiologyRequestBody := createRequestBody(digitalProfileResponse)
	if visiologyRequestBody == nil {
		log.Error("Не удалось создать тело запроса, нет валидных организаций для отправки.")
		return
	}

	jsonBody, err := json.MarshalIndent(visiologyRequestBody, "", " ")
	if err != nil {
		log.Error("Ошибка при маршалировании JSON тела запроса", "error", err)
		return
	}

	// Проверяем переменную окружения DEBUG
	if os.Getenv("DEBUG") == "True" {
		// Если DEBUG=True, выводим тело запроса и не отправляем его
		fmt.Println("Тело запроса (тестовый режим):", string(jsonBody))
		log.Debug("Сформированное тело запроса для Visiology (тестовый режим): ", "jsonBody", string(jsonBody))
		return // Возвращаемся и не отправляем запрос
	}

	// Запрашиваем подтверждение перед отправкой
	fmt.Print("Приступить к обновлению данных Visiology? (Да/Нет): ")
	var userResponse string
	_, _ = fmt.Scanln(&userResponse) //nolint:errcheck
	// Проверяем ввод пользователя
	if strings.ToLower(userResponse) != "да" {
		log.Info("Обновление данных Visiology отменено пользователем.")
		return // Прерываем выполнение, если пользователь не подтвердил
	}

	// В противном случае, отправляем запрос
	response, err := sendRequest(visiologyURL, visiologyAPIVersion, visiologyBearer, jsonBody) //nolint:bodyclose
	if err != nil {
		log.Error("Ошибка при отправке HTTP-запроса", "error: ", err)
		return
	}

	defer closeResponse(response.Body)

	if response.StatusCode != http.StatusOK {
		handleNonOkResponse(response)
	} else {
		log.Info("Данные успешно отправлены на Visiology")
	}
}

// createRequestBody формирует тело запроса на основе ответа от API цифрового профиля.
func createRequestBody(response digitalprofile.GetResponse) []map[string]any {

	for rownum, orgID := range OrgIDs {
		if rownum >= maxIterations {
			break
		}

		for _, org := range response.Organizations {
			if org.ID == orgID {
				// Генерируем данные по всем полям для данной организации
				for _, field := range fields {
					rowData := map[string]any{
						"rownum": rownum,
						"values": []map[string]any{{
							"column": field,
							"value":  org.GetValueByField()[field],
						}},
					}
					requestBody = append(requestBody, rowData)
				}
				break // Прерывание после нахождения нужной организации
			}
		}
	}

	if len(requestBody) == 0 {
		return nil // Нет данных для отправки
	}
	return requestBody
}

// sendRequest отправляет HTTP POST-запрос на указанный URL и возвращает ответ.
func sendRequest(visiologyURL, visiologyAPIVersion, visiologyBearer string, jsonBody []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", visiologyURL+"viqube/databases/DB/tables/KHV_SPO/records/update", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", visiologyBearer)
	req.Header.Add("X-Api-Version", visiologyAPIVersion)

	client := &http.Client{}
	return client.Do(req) //nolint:bodyclose
}

// closeResponse закрывает тело ответа и логирует ошибку, если она произошла.
func closeResponse(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Error("Ошибка закрытия тела ответа", "error", err)
	}
}

// handleNonOkResponse обрабатывает некорректный статус HTTP-ответа.
func handleNonOkResponse(resp *http.Response) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка при чтении тела ответа", "error", err)
	}

	log.Error("Некорректный статус HTTP", "status", resp.StatusCode, "body", string(bodyBytes))
}
