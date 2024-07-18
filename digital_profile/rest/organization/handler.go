package organization

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"visiologyDataUpdate/digital_profile/structs"
)

type Response struct {
	Count         int                    `json:"count"`
	Next          any                    `json:"next"`
	Previous      any                    `json:"previous"`
	Organizations []structs.Organization `json:"results"`
}

// Handler
// Внутри функции мы получаем JSON response из АПИ ЦП по выгрузке организаций
func Handler(digitalProfileUrl string, digitalProfileBearer string) Response {
	req, err := http.NewRequest("GET", digitalProfileUrl+"organizations", nil)
	if err != nil {
		log.Fatal("Ошибка: %v", err)
	}
	req.Header.Add("Authorization", digitalProfileBearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка в ответе.\n[ERROR] -", err)
		panic(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-ok HTTP status:", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка во время обработки JSON:", err)
		panic(err.Error())
	}

	var response Response
	json.Unmarshal(body, &response)

	return response
}
