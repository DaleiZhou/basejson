package basejson

const (
	EOI            = 0x1A
	NOT_MATCH      = -1
	NOT_MATCH_NAME = -2
	UNKNOWN         = 0
	OBJECT         = 1
	ARRAY          = 2
	VALUE          = 3
	END            = 4
	VALUE_NULL     = 5

	UINT32_MAX = uint32(4294967295)
	UINT32_MIN = uint32(0)

	INT32_MAX = int32(2147483647)
	INT32_MIN = int32(-2147483648)

	UINT64_MAX = uint64(18446744073709551615)
	UINT64_MIN = uint64(0)

	INT64_MAX = int64(9223372036854775807)
	INT64_MIN = int64(-9223372036854775808)
)


