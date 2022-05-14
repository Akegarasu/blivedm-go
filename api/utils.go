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
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}
