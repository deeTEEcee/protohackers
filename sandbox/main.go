package main

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

func main() {
	//reqMap := make(map[string]interface{})
	//reqMap["method"] = "isPrime"
	//reqMap["number"] = 123
	var reqObj = Request{"isPrime", 123}
	var request Request
	var err error
	reqBytes, err := json.Marshal(reqObj)
	if err != nil {
		print(err)
	}
	err = json.Unmarshal(reqBytes, &request)
	if err != nil {
		print(err)
	}
	fmt.Println(request)
	fmt.Println(err)
}
