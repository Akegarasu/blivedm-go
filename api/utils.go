package api

import (
	"encoding/json"
	"net/http"
)

func GetJson(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}

func GetJsonWithHeader(url string, headers *http.Header, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	req.Header = *headers
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}
