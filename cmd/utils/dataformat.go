package utils

import "encoding/json"

////////////////////////// Exported funcs //////////////////////////

func ToJson(x interface{}) (string, error) {
	var b []byte
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}
