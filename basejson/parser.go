package basejson

import (
	"sync"
	"errors"
	"fmt"
	"bytes"
)

type Parser interface {
	ParseJson()
}

type JsonParser struct {
	mtx sync.Mutex

	pos int
	len int

	ch byte
	np int

	buf *bytes.Buffer

	token int
	json string
}

func NewJsonParser(str string) *JsonParser {
	_parser := &JsonParser {
		pos: 0,
		len: len(str),
		json: str,
		buf: bytes.NewBuffer([]byte{}),
	}

	_parser.skipWhiteSpace()
	_parser.currentToken()
	return _parser
}

func (this *JsonParser) ParseJSONObject() (*JSONObject, error) {
	if this.currentChar() != '{' {
		return nil, errors.New(fmt.Sprintf("except { at begin of object, pos: %d", this.pos))
	}

	this.next() // skip {
	obj := NewJSONObject()
	for {
		this.skipWhiteSpace()
		ch := this.currentChar()
		var objKey string
		var objValue interface{}

		switch ch {
		case ',': {
			this.next()
			this.skipWhiteSpace()
			continue
		}
		case '"': {
			this.buf.Reset()
			key, err := this.readString()
			if err != nil {
				return nil, err
			}

			this.skipWhiteSpace()
			if this.currentChar() != ':' {
				return nil, errors.New(fmt.Sprintf("Expect : at %d, key = %s", &this.pos, key))
			}
			objKey = key
		}
		case '}': {
			this.buf.Reset()
			this.next()
			return obj, nil
		}
		default: {
			return nil, errors.New(fmt.Sprintf("Expect : at %d, char, %s", this.pos, string(ch)))
		}
		}

		this.next()
		this.skipWhiteSpace()
		ch = this.currentChar()
		println(string(this.currentChar()))

		switch ch {
		case '"': {
			this.buf.Reset()
			value, err := this.readString()
			if err != nil {
				return nil, err
			}
			objValue = value
		}

		case '{': {
			value, err := this.ParseJSONObject()
			if err != nil {
				return nil, err
			}
			objValue = value
		}
		case '[': {
			this.ParseJSONArray()
		}

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-': {
			this.parseNumber()
			switch this.token {
			case LITERAL_INT, LITERAL_LONG: {
				value, err := this.readInt()
				if err != nil {
					return nil, err
				}
				objValue = value
			}
			case LITERAL_FLOAT : {
				value, err := this.readFloat()
				if err != nil {
					return nil, err
				}
				objValue = value
			}
			case LITERAL_DOUBLE: {
				value, err := this.readDouble()
				if err != nil {
					return nil, err
				}
				objValue = value
			}
			}
		}
		} // end of switch
		obj._inner_obj[objKey] = objValue
	}

	return obj, nil
}

func (this *JsonParser) ParseJSONArray() (*JSONArray, error) {
	this.skipWhiteSpace()
	if this.currentChar() != '[' {
		return nil, errors.New(fmt.Sprintf("except [ at the begin of array, pos: %d", this.pos))
	}
	this.next() //skip begin [

	array := NewJsonArray()
	for {
		ch := this.currentChar()

		switch ch {
		case '{': {
			obj, err := this.ParseJSONObject()
			if err != nil {
				return nil, err
			}
			array.Put(obj)
		}
		case ',': {
			this.next()
			this.skipWhiteSpace()
			continue
		}
		}
	}

	return array, nil
}


func (this *JsonParser) Parse() {
	if this.token == LBRACE {
		obj , err := this.ParseJSONObject()
		if err != nil {
			fmt.Println(err.Error())
		} else if str, err := obj.MarshalJSON(); err == nil {
			fmt.Println(string(str))
		}
	} else if this.token == LBRACKET {
		arr, err := this.ParseJSONArray()
		if err != nil {
			fmt.Println(err.Error())
		} else if str, err := arr.MarshalJSON(); err == nil {
			fmt.Println(string(str))
		}
	}
}

func (this *JsonParser) parseNumber() {
	this.np = this.pos
	this.token = LITERAL_INT

	if this.getCharAt(this.np) == '-' {
		this.np ++
	}

	for {
		ch := this.getCharAt(this.np)
		if ch < '0' || ch > '9' {
			break
		}
		this.np ++
	}

	if this.getCharAt(this.np) == '.' {
		this.token = LITERAL_FLOAT
		this.np ++

		for {
			ch := this.getCharAt(this.np)
			if ch < '0' || ch > '9' {
				break
			}
			this.np ++
		}
	}

	ch := this.getCharAt(this.np)
	if ch == 'L' || ch == 'l' {
		this.token = LITERAL_LONG
		this.np ++
	} else if ch == 'D' || ch == 'd' {
		this.token = LITERAL_DOUBLE
		this.np ++
	} else if ch == 'E' || ch == 'e' {
		this.np ++
		ch = this.getCharAt(this.np)
		if ch == '+' || ch == '-' {
			this.np ++
		}

		for {
			ch = this.getCharAt(this.np)
			if ch < '0' || ch > '9' {
				break
			}
			this.np ++
		}

		if ch == 'D' || ch == 'd' {
			this.token = LITERAL_DOUBLE
			this.np ++
		}
	}
}

func (this *JsonParser) currentToken() int {
	tkn := byteToToken(this.currentChar())
	this.token = tkn
	return tkn
}

func (this *JsonParser) nextToken() int {
	tkn := byteToToken(this.next())
	this.token = tkn
	return tkn
}

func (this *JsonParser) readLiteral(literalValue []byte, literalType int) (int, error) {
	literalLen := len(literalValue)
	if literalLen + this.pos >= len(this.json) {
		return 0, errors.New("invalid json string")
	}
	for i := 0; i < literalLen; i++ {
		if this.json[this.pos + i] != literalValue[i] {
			return 0, errors.New("invalid json string")
		}
	}
	this.pos += literalLen
	return literalType, nil
}

func (this *JsonParser) readLong() (int64, error) {
	value := int64(0)
	isNegitive := false
	if this.currentChar() == '-' {
		isNegitive = true
		this.next()
	}

	for {
		ch := this.currentChar()
		if ch == 'L' || ch == 'l' {
			this.next()
			break
		} else if ch >= '0' && ch <= '9' {
			digit := int64(ch - '0')
			if INT64_MAX - digit < value {
				return 0, errors.New(fmt.Sprintf("Can not parse as long value, pos: %d", this.pos))
			}
			value = value * 10 + digit
		} else {
			break
		}
		this.next()
	}

	if isNegitive {
		value = -1 * value
	}

	return value, nil
}

func (this *JsonParser) readInt() (interface{}, error) {
	value64 := uint64(0)

	isNegitive := false
	if this.currentChar() == '-' {
		isNegitive = true
		this.next()
	}

	for {
		ch := this.currentChar()
		println(string(ch))

		if ch == 'L' || ch == 'l' {
			this.next()
			break
		} else if ch >= '0' && ch <= '9' {
			digit := uint64(ch - '0')
			value64 = value64 * 10 + digit
			if (this.token == LITERAL_LONG && uint64(INT64_MAX) < value64) || (this.token == LITERAL_INT && uint64(INT32_MAX) < value64) {
				return 0, errors.New(fmt.Sprintf("Parse value got error: overflows int32(int64), pos: %d", this.pos))
			}

		} else {
			break
		}
		this.next()
	}

	returnValue := int64(value64)
	if isNegitive {
		returnValue = -1 * returnValue
	}

	if this.token == LITERAL_INT {
		return int32(returnValue), nil
	}
	return returnValue, nil
}

func (this *JsonParser) readFloat() (float32, error) {

	return 0, nil
}

func (this *JsonParser) readDouble() (float64, error) {
	return 0, nil
}

func (this *JsonParser) readString() (string, error) {
	var ch byte
	for {
		ch = this.next()
		if ch == '"' {
			break
		}

		if ch == EOI {
			if !this.isEOF() {
				this.putChar(EOI)
				continue
			}
			return "", errors.New(fmt.Sprintf("unclosed string : %c", ch))
		}

		if ch == '\\' {
			this.putChar('\\')
			ch = this.next()
			this.putChar(ch)
		}
		this.putChar(ch)
	}
	this.token = LITERAL_STRING
	this.ch = this.next()
	return this.buf.String(), nil
}

func (this *JsonParser) putChar(ch byte) {
	this.buf.WriteByte(ch)
}

func (this *JsonParser) next() byte {
	this.pos ++
	return this.getCharAt(this.pos)
}

func (this *JsonParser) currentChar() byte {
	return this.getCharAt(this.pos)
}

func (this *JsonParser) getCharAt(pos int) byte {
	if pos >= this.len {
		return EOI
	}
	return this.json[pos]
}

func (this *JsonParser) isEOF() bool {
	return this.pos == this.len || this.ch == EOI && this.pos + 1 == this.len
}

func isWhiteSpace(ch byte) bool {
	return  ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

func (this *JsonParser) skipWhiteSpace() {
	for {
		if !isWhiteSpace(this.currentChar()) {
			break
		}
		this.next()
	}
}

func (this *JsonParser) parseNull() {
	this.token = NULL
}
