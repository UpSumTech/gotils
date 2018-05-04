package utils

import (
	"encoding/json"
	"io/ioutil"
)

////////////////////////// Exported funcs //////////////////////////

func ToJson(x interface{}) (string, error) {
	var b []byte
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}

func IsValidJson(f string) bool {
	j, err := ioutil.ReadFile(f)
	if err != nil {
		return false
	}
	var js json.RawMessage
	return json.Unmarshal([]byte(j), &js) == nil
}
