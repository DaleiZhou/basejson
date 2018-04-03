package basejson

import "encoding/json"

type JSONObject struct {
	_inner_obj map[string]interface{}
}

func NewJSONObject() (*JSONObject) {
	return &JSONObject{
		_inner_obj: make(map[string]interface{}),
	}
}

func(obj *JSONObject) Size() int {
	return len(obj._inner_obj)
}

func(obj *JSONObject) IsEmpty() bool {
	return obj.Size() == 0
}

func(obj *JSONObject) ContainsKey(key string) bool {
	_, ok := obj._inner_obj[key]
	return ok
}

func(obj *JSONObject) Get(key string) interface{} {
	if val, ok := obj._inner_obj[key]; ok {
		return val
	}
	return nil
}

func(obj *JSONObject) GetJSONObject(key string) *JSONObject {
	if val, ok := obj._inner_obj[key]; ok {
		if ret, ok := val.(*JSONObject); ok {
			return ret
		}
		return nil
	}
	return nil
}

func(obj *JSONObject) GetJSONArray(key string) *JSONArray {
	if val, ok := obj._inner_obj[key]; ok {
		if ret, ok := val.(*JSONArray); ok {
			return ret
		}
		return nil
	}
	return nil
}

func(obj *JSONObject) Put(key string, value interface{}) {
	obj._inner_obj[key] = value
}

func(obj *JSONObject) PutAll(m map[string]interface{}) {
	for key, val := range m {
		obj._inner_obj[key] = val
	}
}

func(obj *JSONObject) Delete(key string) {
	delete(obj._inner_obj, key)
}

func(obj *JSONObject) DeleteAll(keys []string) {
	for _, key := range keys {
		obj.Delete(key)
	}
}

func(obj *JSONObject) Clear() {
	obj._inner_obj =  make(map[string]interface{})
}

func(obj *JSONObject) MarshalJSON() ([]byte, error) {
	return json.Marshal(&obj._inner_obj)
}

