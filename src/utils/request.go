package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"reflect"
)

func MakeRequest(method, url string, headers map[string]string, body interface{}) (map[string]interface{}, error) {
	client := resty.New()
	req := client.R().
		SetHeaders(headers).
		SetBody(body)

	resp, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	defer resp.RawResponse.Body.Close()
	fmt.Println("Body type", reflect.TypeOf(resp.Body))

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
