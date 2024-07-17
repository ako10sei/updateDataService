package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Count         int            `json:"count"`
	Next          any            `json:"next"`
	Previous      any            `json:"previous"`
	Organizations []Organization `json:"results"`
}

func main() {
	url := "https://xn--n1abf.xn--33-6kcadhwnl3cfdx.xn--p1ai/digital_profile/api/v1.0.0/organizations"
	var bearer = "Bearer " + "Ed3SlmOQvZldnYDL7aE4qDX1lBk0Eo"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка в ответе.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка во время обработки JSON:", err)
	}

	var response Response

	look := json.Unmarshal(body, &response)

	if look != nil {
		fmt.Println(look)
	}

	for i, p := range response.Organizations {
		fmt.Println("Organization", i+1, ":", p.Name, p.ShortName)
	}
}
