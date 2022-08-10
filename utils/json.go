package utils

import "encoding/json"

func UnmarshalStr(str string, v interface{}) error {
	return json.Unmarshal(StringToBytes(str), v)
}
