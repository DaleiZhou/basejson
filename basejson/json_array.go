package basejson

import "encoding/json"

type JSONArray struct {
	_inner_array []interface{}
}

func NewJsonArray() *JSONArray {
	return &JSONArray{
		_inner_array: make([]interface{}, 0),
	}
}

func(arr *JSONArray) Size() int {
	return len(arr._inner_array)
}

func(arr *JSONArray) IsEmpty() bool {
	return len(arr._inner_array) == 0
}

func(arr *JSONArray) Put(obj interface{}) {
	arr._inner_array = append(arr._inner_array, obj)
}

func(arr *JSONArray) PutJSONObject(obj *JSONObject) {
	arr._inner_array = append(arr._inner_array, obj)
}

func(arr *JSONArray) Get(index int) *interface{} {
	obj := arr._inner_array[index]
	return &obj
}

func(arr *JSONArray) GetJSONObject(index int) *JSONObject {
	obj := arr._inner_array[index]
	if v, ok := obj.(JSONObject); ok {
		return &v
	}
	return nil
}

func(arr *JSONArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(&arr._inner_array)
}




