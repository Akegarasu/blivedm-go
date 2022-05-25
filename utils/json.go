package utils

import "encoding/json"

func UnmarshalString(str string, v interface{}) error {
	return json.Unmarshal(StringToBytes(str), v)
}
