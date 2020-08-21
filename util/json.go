package util

import (
	"bytes"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func ParseJson(data []byte) map[string]interface{} {
	var result map[string]interface{}
	d := JSON.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	d.Decode(&result)
	return result
}

func ToJson(object interface{}) string {
	json, err := JSON.Marshal(object)
	if err != nil {
		fmt.Println("ToJson Error：", err)
		return "{}"
	}
	return string(json)
}
