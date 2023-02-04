package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func Contains[K comparable](s []K, item K) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}

func MoveToFirstPositionOfSlice[K comparable](slice []K, item K) []K {
	if len(slice) == 0 || (slice)[0] == item {
		return slice
	}
	if (slice)[len(slice)-1] == item {
		slice = append([]K{item}, (slice)[:len(slice)-1]...)
		return slice
	}
	for p, x := range slice {
		if x == item {
			slice = append([]K{item}, append((slice)[:p], (slice)[p+1:]...)...)
			break
		}
	}
	return slice
}

func ifThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func Convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = Convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = Convert(v)
		}
	}
	return i
}

func ObjectFronJSON(jsonString string, result interface{}) error {
	return json.Unmarshal([]byte(jsonString), &result)
}

func Sha256String(name string) string {
	hash := sha256.Sum256([]byte(name))
	return hex.EncodeToString(hash[:])
}

// Given an input slice of bytes representing an arbitrary JSON object and a slice of strings containing keys
// which should not exist in the input JSON object, remove these keys from original object.
func removeKeysFromJSONObject(input *map[string]json.RawMessage, keys []string) {
	for _, key := range keys {
		delete(*input, key)
	}
}

func RemoveKeysFromJSONObjectBytes(input *[]byte, keys []string) error {
	var output map[string]json.RawMessage
	if err := json.Unmarshal(*input, &output); err != nil {
		return err
	}
	err := RemoveKeysFromJSONObject(&output, keys)
	if err != nil {
		return err
	}
	outputBytes, err := json.Marshal(&output)
	if err != nil {
		return err
	}
	*input = outputBytes
	return nil
}

func RemoveKeysFromJSONObject(input *map[string]json.RawMessage, keys []string) error {
	removeKeysFromJSONObject(input, keys)
	return nil
}

func mapFuncToJsonObjectArray(fn func(input *map[string]json.RawMessage) error, jsonArray *[]map[string]json.RawMessage) error {
	for _, jsonArrayItem := range *jsonArray {
		err := fn(&jsonArrayItem)
		if err != nil {
			return err
		}
	}
	return nil
}

// Given a function which accepts an arbitrary JSON object, map this function and its outputs onto a provided
// arbitrary array of JSON objects.
func MapFuncToJsonObjectArrayBytes(fn func(input *map[string]json.RawMessage) error, jsonArray *[]byte) error {
	var output []map[string]json.RawMessage
	if err := json.Unmarshal(*jsonArray, &output); err != nil {
		return err
	}
	err := MapFuncToJsonObjectArray(fn, &output)
	if err != nil {
		return err
	}
	// https://stackoverflow.com/a/24229303/4562156
	// The value passed to json.Marshal must be a pointer for json.RawMessage to work properly.
	outputBytes, err := json.Marshal(&output)
	if err != nil {
		return err
	}
	*jsonArray = outputBytes
	return nil
}

func MapFuncToJsonObjectArray(fn func(input *map[string]json.RawMessage) error, jsonArray *[]map[string]json.RawMessage) error {
	return mapFuncToJsonObjectArray(fn, jsonArray)
}
