package security

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CveResponse struct {
	Cvss float64 `json:"cvss"`
}

func SearchDetailCve(cve string) (CveResponse, error) {

	var response CveResponse

	url := fmt.Sprintf("https://cve.circl.lu/api/cve/%s", cve)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "PostmanRuntime/7.29.0")

	resp, err := client.Do(req)

	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return response, err
	}

	return response, nil

}
