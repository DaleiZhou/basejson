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

	lexer
}

type lexer struct {
	pos int
	len int

	ch byte
	np int

	buf *bytes.Buffer

	token int
	json string
}

func NewJsonParser(str string) *jsonParser {
	_parser := &jsonParser{
		lexer: lexer {
			pos: 0,
			len: len(str),
			json: str,
			np : -1,
			buf: bytes.NewBuffer([]byte{}),
		},
	}

	_parser.skipWhiteSpace()
	_parser.currentToken()
	return _parser
}

func (this *lexer) ParseJSONObject() (*JSONObject, error) {
	if this.currentChar() != '{' {
		return nil, errors.New(fmt.Sprintf("except { at begin of object, pos: %d", this.pos))
	}

	this.next() // skip begin '{'
	obj := NewJSONObject()
	for {
		this.skipWhiteSpace()

		var objKey string
		var objValue interface{}

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
				return nil, errors.New(fmt.Sprintf("Expect : at pos: %d, key = %s", this.pos, key))
			}
			objKey = key
		}
		case RBRACE: {
			this.buf.Reset()
			this.next()
			return obj, nil
		}
		default: {
			return nil, errors.New(fmt.Sprintf("Error char: %s at pos: %d, %s", string(this.currentChar()), this.pos, this.json[this.pos:]))
		}
		}

		this.next()
		this.skipWhiteSpace()

		switch this.currentToken() {
		case TRUE:{
			err := this.readLiteral("true")
			if err != nil {
				return nil, err
			}
			objValue = true
		}
		case FALSE: {
			err := this.readLiteral("false")
			if err != nil {
				return nil, err
			}
			objValue = false
		}
		case NULL :{
			err := this.readLiteral("null")
			if err != nil {
				return nil, err
			}
			objValue = nil
		}
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
			value, err := this.ParseJSONArray()
			if err != nil {
				return nil, err
			}
			objValue = value
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

func (this *lexer) ParseJSONArray() (*JSONArray, error) {
	this.skipWhiteSpace()
	if this.currentChar() != '[' {
		return nil, errors.New(fmt.Sprintf("except [ at the begin of array, pos: %d", this.pos))
	}
	this.next() //skip begin [

	array := NewJsonArray()
	for {
		this.skipWhiteSpace()
		token := this.currentToken()
		switch token {
		case TRUE:{
			err := this.readLiteral("true")
			if err != nil {
				return nil, err
			}
			array.Put(true)
		}
		case FALSE: {
			err := this.readLiteral("false")
			if err != nil {
				return nil, err
			}
			array.Put(false)
		}
		case NULL :{
			err := this.readLiteral("null")
			if err != nil {
				return nil, err
			}
			array.Put(nil)
		}
		case LBRACKET: {
			obj, err := this.ParseJSONArray()
			if err != nil {
				return nil, err
			}
			array.Put(obj)
		}
		case COMMA : {
			this.next()
			this.skipWhiteSpace()
			continue
		}
		case LBRACE: {
			obj, err := this.ParseJSONObject()
			if err != nil {
				return nil, err
			}
			array.Put(obj)
		}
		case RBRACKET: {
			this.buf.Reset()
			this.next()
			return array, nil
		}
		case DOUBLE_QUOTES: {
			this.buf.Reset()
			strValue, err := this.readString()
			if err != nil {
				return nil, err
			}
			array.Put(strValue)
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

				array.Put(intValue)
			}
			case LITERAL_FLOAT, LITERAL_DOUBLE: {
				floatVal, err := strconv.ParseFloat(string(this.json[this.pos: this.np]), 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("parse number got error, pos: %d", this.pos))
				}

				array.Put(floatVal)
			}
			}

			this.pos = this.np
			if this.token == LITERAL_LONG {
				this.next()
			}
			this.np = 0
		}
		}
	}

	return array, nil
}


func (this *lexer) Parse() {
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

func (this *lexer) scanNumberToken() {
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

func (this *lexer) currentToken() int {
	tkn := byteToToken(this.currentChar())
	this.token = tkn
	return tkn
}

func (this *lexer) nextToken() int {
	tkn := byteToToken(this.next())
	this.token = tkn
	return tkn
}

func (this *lexer) readLiteral(literalValue string) error {
	literalLen := len(literalValue)
	for idx := 0; idx < literalLen; idx++ {
		ch := this.currentChar()
		if ch != literalValue[idx] {
			return errors.New(fmt.Sprintf("expect literal value: %s, at pos: %d", string(literalValue), this.pos))
		}
		this.next()
	}
	return nil
}

func (this *lexer) readString() (string, error) {
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

func (this *lexer) putChar(ch byte) {
	this.buf.WriteByte(ch)
}

func (this *lexer) next() byte {
	this.pos ++
	return this.getCharAt(this.pos)
}

func (this *lexer) currentChar() byte {
	return this.getCharAt(this.pos)
}

func (this *lexer) getCharAt(pos int) byte {
	if pos >= this.len {
		return EOI
	}
	return this.json[pos]
}

func (this *lexer) isEOF() bool {
	return this.pos == this.len || this.ch == EOI && this.pos + 1 == this.len
}

func isWhiteSpace(ch byte) bool {
	return  ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}

func (this *lexer) skipWhiteSpace() {
	for {
		if !isWhiteSpace(this.currentChar()) {
			break
		}
		this.next()
	}
}

func (this *lexer) parseNull() {
	this.token = NULL
}