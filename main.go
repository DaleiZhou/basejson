package main


import (
	"baseJson/basejson"
)

func main(){
	//str1 := "{\"key\"    : 23423423434234, \"key1\"    : \"value2\", \"key2\" :    { \"inner_key1\" : \"inner_value1\"  } , \"key3\"    : \"value3\",}"
	//str2 := "[{\"key\"    : \"value\", , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }]"
	//str3 := "[{\"key\"    : 123123L, , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }, {\"key\"    : \"value\", , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }]"

	//str4 := "[    \"123123\"   ,    \"sdfsdf123123\"    , \"ge34g34g\"    ]"
	//str5 := "[    123123   ,    2.342342e34    , 234234234    ]"
	//str6 := "[[[\"test\"], [ \"test2\", \"test3\" ]]]"

	str7 := "[ true, false, null]"
	parser := basejson.NewJsonParser(str7)
	parser.Parse()
}