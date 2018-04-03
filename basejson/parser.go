package basejson

import (
	"sync"
	"errors"
	"fmt"
	"bytes"
	"strconv"
)

type Parser interface {
	ParseJson()
}

type jsonParser struct {
	mtx sync.Mutex

	pos int
	len int

	ch byte
	np int

	buf *bytes.Buffer

	token int
	json string
}

func NewJsonParser(str string) *jsonParser {
	_parser := &jsonParser {
		pos: 0,
		len: len(str),
		json: str,
		buf: bytes.NewBuffer([]byte{}),
	}

	_parser.skipWhiteSpace()
	_parser.currentToken()
	return _parser
}

func (this *jsonParser) ParseJSONObject() (*JSONObject, error) {
	if this.currentChar() != '{' {
		return nil, errors.New(fmt.Sprintf("except { at begin of object, pos: %d", this.pos))
	}

	this.next() // skip begin '{'
	obj := NewJSONObject()
	for {
		this.skipWhiteSpace()

		var objKey string
		var objValue interface{}
		fmt.Println("char: ", string(this.currentChar()), "token:", this.currentToken())

		switch this.currentToken() {
		case COMMA : {
			this.next()
			this.skipWhiteSpace()
			continue
		}
		case DOUBLE_QUOTES: {
			this.buf.Reset()
			key, err := this.readString()
			if err != nil {
				return nil, err
			}

			this.skipWhiteSpace()
			if this.currentToken() != COLON {
				return nil, errors.New(fmt.Sprintf("Expect : at %d, key = %s", &this.pos, key))
			}
			objKey = key
		}
		case RBRACE: {
			this.buf.Reset()
			this.next()
			return obj, nil
		}
		default: {
			return nil, errors.New(fmt.Sprintf("Error char: %s at pos: %d", string(this.currentChar()), this.pos))
		}
		}

		this.next()
		this.skipWhiteSpace()

		switch this.currentToken() {
		case DOUBLE_QUOTES: {
			this.buf.Reset()
			value, err := this.readString()
			if err != nil {
				return nil, err
			}
			objValue = value
		}

		case LBRACE: {
			value, err := this.ParseJSONObject()
			if err != nil {
				return nil, err
			}
			objValue = value
		}
		case LBRACKET: {
			this.ParseJSONArray()
		}

		case DIGIT: {
			this.scanNumberToken()
			if this.np <= this.pos {
				return nil, errors.New(fmt.Sprintf("parse number got error, pos: %d", this.pos))
			}

			switch this.token {
			case LITERAL_INT, LITERAL_LONG: {
				intValue, err := strconv.ParseInt(string(this.json[this.pos: this.np]), 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("parse number got error, pos: %d", this.pos))
				}

				objValue = intValue
			}
			case LITERAL_FLOAT, LITERAL_DOUBLE: {
				floatVal, err := strconv.ParseFloat(string(this.json[this.pos: this.np]), 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("parse number got error, pos: %d", this.pos))
				}

				objValue = floatVal
			}
			}

			this.pos = this.np
			if this.token == LITERAL_LONG {
				this.pos ++
			}
			this.np = 0
		}
		} // end of switch
		obj._inner_obj[objKey] = objValue
	}

	return obj, nil
}

func (this *jsonParser) ParseJSONArray() (*JSONArray, error) {
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


func (this *jsonParser) Parse() {
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

func (this *jsonParser) scanNumberToken() {
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
	}
}

func (this *jsonParser) currentToken() int {
	tkn := byteToToken(this.currentChar())
	this.token = tkn
	return tkn
}

func (this *jsonParser) nextToken() int {
	tkn := byteToToken(this.next())
	this.token = tkn
	return tkn
}

func (this *jsonParser) readLiteral(literalValue []byte, literalType int) (int, error) {
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

/**
 * Deprecated. Using strconv.ParseInt instead
 */
func (this *jsonParser) readInt() (interface{}, error) {
	value64 := uint64(0)

	isNegitive := false
	if this.currentChar() == '-' {
		isNegitive = true
		this.next()
	}

	for {
		ch := this.currentChar()
		if ch < '0' || ch > '9' {
			if ch == 'L' || ch == 'l' {
				this.next()
			}
			break
		}

		digit := uint64(ch - '0')
		value64 = value64 * 10 + digit
		if (this.token == LITERAL_LONG && uint64(INT64_MAX) < value64) || (this.token == LITERAL_INT && uint64(INT32_MAX) < value64) {
			return 0, errors.New(fmt.Sprintf("Parse value got error: overflows int32(int64), pos: %d", this.pos))
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

func (this *jsonParser) readFloat() (float32, error) {

	return 0, nil
}

func (this *jsonParser) readDouble() (float64, error) {
	return 0, nil
}

func (this *jsonParser) readString() (string, error) {
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

func (this *jsonParser) putChar(ch byte) {
	this.buf.WriteByte(ch)
}

func (this *jsonParser) next() byte {
	this.pos ++
	return this.getCharAt(this.pos)
}

func (this *jsonParser) currentChar() byte {
	return this.getCharAt(this.pos)
}

func (this *jsonParser) getCharAt(pos int) byte {
	if pos >= this.len {
		return EOI
	}
	return this.json[pos]
}

func (this *jsonParser) isEOF() bool {
	return this.pos == this.len || this.ch == EOI && this.pos + 1 == this.len
}

func isWhiteSpace(ch byte) bool {
	return  ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

func (this *jsonParser) skipWhiteSpace() {
	for {
		if !isWhiteSpace(this.currentChar()) {
			break
		}
		this.next()
	}
}

func (this *jsonParser) parseNull() {
	this.token = NULL
}
