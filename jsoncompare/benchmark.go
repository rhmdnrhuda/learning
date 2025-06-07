package jsoncompare

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
)

// Standard library JSON unmarshal function
func StandardJsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// JsonIter unmarshal function
func JsonIterUnmarshal(data []byte, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, v)
}
