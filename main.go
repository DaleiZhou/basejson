package main


import (
	"baseJson/basejson"
)


func main(){
	str1 := "{\"key\"    : -0.1L, \"key1\"    : \"value2\", \"key2\" :    { \"inner_key1\" : \"inner_value1\"  } }"
	//str2 := "[{\"key\"    : \"value\", , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }]"
	//str3 := "[{\"key\"    : 123123L, , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }, {\"key\"    : \"value\", , , , \"key1\"    : \"value2\", \"key2\" : { \"inner_key1\" : \"inner_value1\"  } }]"

	parser := basejson.NewJsonParser(str1)
	parser.Parse()
}