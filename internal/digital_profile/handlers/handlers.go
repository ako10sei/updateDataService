package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"visiologyDataUpdate/internal/digital_profile/structs"
	"visiologyDataUpdate/internal/log"
)

// GetResponse представляет ответ от API цифрового профиля.
type GetResponse struct {
	Count         int                    `json:"count"`
	Next          any                    `json:"next"`
	Previous      any                    `json:"previous"`
	Organizations []structs.Organization `json:"results"`
}

// GetHandler отправляет GET-запрос на указанный URL с указанным маркером доступа,
// обрабатывает ответ и возвращает структуру GetResponse, содержащую данные организаций.
func GetHandler(digitalProfileURL, digitalProfileBearer string) GetResponse {
	log.Info("Отправка GET-запроса на API цифрового профиля", "url: ", digitalProfileURL+"organizations")

	req, err := createRequest(digitalProfileURL+"organizations", digitalProfileBearer)
	if err != nil {
		log.Error("Ошибка создания HTTP-запроса", "error: ", err)
		return GetResponse{}
	}

	resp, err := sendRequest(req) //nolint:bodyclose
	if err != nil {
		return GetResponse{}
	}

	defer closeResponse(resp.Body)

	if resp.StatusCode != http.StatusOK {
		handleNonOKResponse(resp)
		return GetResponse{}
	}

	return parseResponse(resp.Body)
}

func createRequest(url, bearer string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", bearer)
	return req, nil
}

func sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req) //nolint:bodyclose
	return resp, err
}

func handleNonOKResponse(resp *http.Response) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Ошибка во время чтения тела ответа", "error: ", err)
	}
	log.Error("Некорректный статус HTTP", "status: ", resp.StatusCode, "body: ", string(bodyBytes))
}

func parseResponse(body io.ReadCloser) GetResponse {
	data, err := io.ReadAll(body)
	if err != nil {
		log.Error("Ошибка во время чтения тела ответа", "error: ", err)
		return GetResponse{}
	}

	var response GetResponse
	if err := json.Unmarshal(data, &response); err != nil {
		log.Error("Ошибка десериализации ответа", "error: ", err, "body: ", string(data))
		return GetResponse{}
	}

	log.Info("Ответ получен успешно", "count: ", response.Count)
	return response
}

func closeResponse(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Error("Ошибка закрытия тела ответа", "error: ", err)
	}
}
