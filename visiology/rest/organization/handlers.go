package organization

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	digitalprofile "visiologyDataUpdate/digital_profile/rest/organization"
	"visiologyDataUpdate/visiology/structs"
)

type GetResponse struct {
	Columns []structs.Column `json:"columns"`
	Values  [][]any          `json:"values"`
}

// GetHandler
// Внутри функции мы получаем JSON response из АПИ Visiology и возвращаем его в виде структуры GetResponse
func GetHandler(visiologyUrl string, visiologyBearer string, visiologyApiVersion string) GetResponse {
	var response GetResponse
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
		fmt.Println("GetResponse body:", string(bodyBytes))
		return GetResponse{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка во время обработки JSON:", err)
		panic(err.Error())
	}
	json.Unmarshal(body, &response)

	return response
}

func PostHandler(
	digitalProfileResponse digitalprofile.GetResponse,
	visiologyResponse GetResponse,
	visiologyUrl string,
	visiologyApiVersion string,
	visiologyBearer string) {

	requestBody := []map[string]interface{}{}
	for i, d := range digitalProfileResponse.Organizations {
		for _, col := range visiologyResponse.Columns {
			rowData := map[string]interface{}{
				"rownum": i,
				"values": []map[string]interface{}{
					{
						"column": col.CoollegeId,
						"value":  d.ID,
					},
				},
			}
			requestBody = append(requestBody, rowData)
		}
	}

	jsonBody, err := json.MarshalIndent(requestBody[:20], "", "  ")
	if err != nil {
		return
	}

	fmt.Println(string(jsonBody))
	//jsonBody, err := json.Marshal(requestBody)
	//if err != nil {
	//	return
	//}
	//
	//req, err := http.NewRequest("POST", visiologyUrl, bytes.NewBuffer(jsonBody))
	//if err != nil {
	//	return
	//}
	//
	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer "+visiologyBearer)
	//req.Header.Set("X-Visiology-Api-Version", visiologyApiVersion)
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	return
	//}
	//defer resp.Body.Close()
}
