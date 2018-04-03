package basejson

const (
	ERROR                = 1;
	LITERAL_INT          = 2;
	LITERAL_FLOAT        = 3;


	LITERAL_STRING       = 4;

	LITERAL_ISO8601_DATE = 5;

	TRUE                 = 6;

	FALSE                = 7;

	NULL                 = 8;

	NEW                  = 9;

	LPAREN               = 10; // ("("),

	RPAREN               = 11; // (")"),

	LBRACE               = 12; // ("{"),

	RBRACE               = 13; // ("}"),

	LBRACKET             = 14; // ("["),

	RBRACKET             = 15; // ("]"),

	COMMA                = 16; // (","),

	COLON                = 17; // (":"),

	IDENTIFIER           = 18;

	FIELD_NAME           = 19;

	EOF                  = 20;
	SET                  = 21;
	TREE_SET             = 22;

	UNDEFINED            = 23; // undefined
	SEMI                 = 24;
	DOT                  = 25;
	HEX                  = 26;

	LITERAL_LONG         = 27;
	LITERAL_DOUBLE       = 28;

	LITERAL_SPACE        = 29; // space
	LITERAL_TAB          = 30; // \t
	LITERAL_RETURN       = 31; // \r
	LITERAL_NEXT         = 32; // \n

	DIGIT                = 33;
	NEGITIVE             = 34;

	DOUBLE_QUOTES        = 35; // "
	SIGNLE_QUOTES        = 36; // '
)

var tokenMap = map[byte]int{
	'(' :  LPAREN,
	')' :  RPAREN,
	'{' :  LBRACE,
	'}' :  RBRACE,
	'[' :  LBRACKET,
	']' :  RBRACKET,
	',' :  COMMA,
	':' :  COLON,
	' ' :  LITERAL_SPACE,
	'\t':  LITERAL_TAB,
	'\r':  LITERAL_RETURN,
	'\n':  LITERAL_NEXT,
	'n' :  NULL,
	'f' :  FALSE,
	't' :  TRUE,
	'0' :  DIGIT,
	'1' :  DIGIT,
	'2' :  DIGIT,
	'3' :  DIGIT,
	'4' :  DIGIT,
	'5' :  DIGIT,
	'6' :  DIGIT,
	'7' :  DIGIT,
	'8' :  DIGIT,
	'9' :  DIGIT,
	'-' :  DIGIT,
	'"' :  DOUBLE_QUOTES,
	'\'':  SIGNLE_QUOTES,
	0x1A:  EOI,
}

func byteToToken(ch byte) int {
	tkn, ok := tokenMap[ch]
	if ok {
		return tkn
	}
	return UNKNOWN
}

func tokenName(value int) string {
	switch value {
	case ERROR:
		return "error"
	case LITERAL_INT:
		return "int"
	case LITERAL_FLOAT:
		return "float"
	case LITERAL_STRING:
		return "string"
	case LITERAL_ISO8601_DATE:
		return "iso8601"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case NULL:
		return "null"
	case NEW:
		return "new"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LBRACE:
		return "{"
	case RBRACE:
		return "}"
	case LBRACKET:
		return "["
	case RBRACKET:
		return "]"
	case COMMA:
		return ","
	case COLON:
		return ":"
	case IDENTIFIER:
		return "ident"
	case FIELD_NAME:
		return "fieldName"
	case EOF:
		return "EOF"
	case SET:
		return "Set"
	case TREE_SET:
		return "TreeSet"
	case UNDEFINED:
		return "undefined"
	case SEMI:
		return ";"
	case DOT:
		return "."
	case HEX:
		return "hex"
	default:
		return "Unknown"
	}
}

