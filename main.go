package main

//https://stackoverflow.com/questions/40343471/how-to-check-if-interface-is-a-slice
//https://www.json.org/json-en.html
//https://datatracker.ietf.org/doc/html/rfc7159
//https://earthly.dev/blog/jq-select/

//value: object, array, number, or string, "false", "null", "true"

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	jsonFile, err := os.Open("file2.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	//var data map[string]any
	var data any
	err = json.Unmarshal([]byte(byteValue), &data)
	if err != nil {
		fmt.Println(err)
	}

	switch d := data.(type) {
	case string:
		fmt.Println("string ", d)
	case []interface{}:
		fmt.Println("is an array:", d)
		for i, u := range d {
			fmt.Println(i, u)
		}
	case map[string]interface{}:
		fmt.Println("is a map")
	default:
		fmt.Println("is of a type I don't know how to handle")
	}

	//if reflect.TypeOf(data).Kind() == reflect.Slice {
	//	for i := range data {
	//		fmt.Printf("%s\n", data[i])
	//
	//

	// reflect.TypeOf(data).Kind() == reflect.Map {
	//	for k, v := range data {
	//		fmt.Printf("%s: %s\n", k, v)
	//	}
	//

	//for k, v := range data {
	//	fmt.Printf("KEY: %s Value: %s\n", k, v)
	//}
}
