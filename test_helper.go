package simple_query

import (
	"encoding/json"
	"reflect"
)

func deepEqual(value1 interface{}, value2 interface{}) bool {
	var (
		val1  interface{}
		val1B []byte
		val2  interface{}
		val2B []byte
	)

	if reflect.DeepEqual(value1, value2) {
		return true
	}

	val1B, _ = json.Marshal(value1)
	json.Unmarshal(val1B, &val1)
	val2B, _ = json.Marshal(value2)
	json.Unmarshal(val2B, &val2)

	return reflect.DeepEqual(val1, val2)
}
