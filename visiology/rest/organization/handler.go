package organization

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"visiologyDataUpdate/visiology/structs"
)

type Response struct {
	Columns []structs.Column `json:"columns"`
	Values  [][]any          `json:"values"`
}

// Handler
// Внутри функции мы получаем JSON response из АПИ ЦП по выгрузке организаций
func Handler(visiologyUrl string, visiologyBearer string, visiologyApiVersion string) Response {
	var response Response
	req, err := http.NewRequest("GET", visiologyUrl, nil)
	if err != nil {
		log.Fatal("Ошибка: %v", err)
	}
	req.Header.Add("Authorization", visiologyBearer)
	req.Header.Add("x-api-version", visiologyApiVersion)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка в ответе.\n[ERROR] -", err)
		panic(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Ошибка во время чтения тела ответа:", err)
			panic(err.Error())
		}
		fmt.Println("Non-ok HTTP status:", resp.StatusCode)
		fmt.Println("Response body:", string(bodyBytes))
		return Response{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка во время обработки JSON:", err)
		panic(err.Error())
	}
	json.Unmarshal(body, &response)

	return response
}
