package utils

import "encoding/json"

////////////////////////// Exported funcs //////////////////////////

func ToJson(x interface{}) string {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		CheckErr(err.Error())
	}
	return string(b)
}
