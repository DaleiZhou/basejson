package basejson

type jsonValue interface {
	MarshalJSON() ([]byte, error)
}

type literalString struct {
	value string
}

func (lv *literalString) MarshalJSON() ([]byte, error) {
	return []byte(lv.value), nil
}

type literalNull struct {
	value interface{}
}

func (lv *literalNull) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

type literalBool struct {
	value bool
}

func (lv *literalBool) MarshalJSON() ([]byte, error) {
	value := "false"
	if lv.value {
		value = "true"
	}
	return []byte(value), nil
}

type literalInt struct {
	value int64
}

//TODO
func (lv *literalInt) MarshalJSON() ([]byte, error) {
	return nil, nil
}

type literalFloat struct {
	value float64
}

//TODO
func (lv *literalFloat) MarshalJSON() ([]byte, error) {
	return nil, nil
}