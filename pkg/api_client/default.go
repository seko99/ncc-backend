package api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorMessage `json:"error"`
}

type DefaultClient struct {
	apiURL   string
	apiToken string
}

func (uc DefaultClient) Do(method string, url string, t interface{}, payloads ...interface{}) error {
	client := &http.Client{}

	var err error
	var b []byte

	if len(payloads) == 1 {
		b, err = json.Marshal(payloads[0])
		if err != nil {
			return err
		}
	}

	body := bytes.NewReader(b)
	req, err := http.NewRequest(method, uc.apiURL+"/"+url, body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+uc.apiToken)
	req.Header.Add("Content-Type", "Application/json")

	response, err := client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		b, err = io.ReadAll(response.Body)

		var errorResponse ErrorResponse
		err = json.Unmarshal(b, &errorResponse)
		if err != nil {
			return fmt.Errorf("can't unmarshal error response: %w body: %s", err, string(b))
		}

		return fmt.Errorf("error response: %d %+v", response.StatusCode, errorResponse)
	}

	err = json.NewDecoder(response.Body).Decode(t)
	if err != nil {
		return fmt.Errorf("can't unmarshal response: %w", err)
	}

	return nil
}
