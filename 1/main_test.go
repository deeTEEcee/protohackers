package main

import (
	"encoding/json"
	"testing"
)

func TestProcessSuccess(t *testing.T) {
	var err error
	reqMap := make(map[string]interface{})
	reqMap["method"] = "isPrime"
	reqMap["number"] = 123
	request, err := json.Marshal(reqMap)
	responseBytes, ok := process(request)
	if !ok {
		t.Errorf("Should've succeeded but failed with %s. Err: %s", responseBytes, err)
	}
	var response Response
	err = json.Unmarshal(responseBytes, &response)
	if response.Method != "isPrime" {
		t.Errorf("Method name '%s' is wrong", response.Method)
	}
	if response.IsPrime {
		t.Errorf("%d is not a prime", 123)
	}
}

func TestProcessFailWithBadJson(t *testing.T) {
	request := make([]byte, 1024)
	response, ok := process(request)
	if ok {
		t.Errorf("Should've failed but succeeded with %s", response)
	}
}

func TestProcessFailWithMissingField(t *testing.T) {
	reqMap := make(map[string]interface{})
	reqMap["method"] = "isPrime"
	request, err := json.Marshal(reqMap)
	response, ok := process(request)
	if ok {
		t.Errorf("Should've failed but succeeded with %s. Err: %s", response, err)
	}
}

func TestIsPrime(t *testing.T) {
	primes := []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73,
		79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173,
		179, 181, 191, 193, 197, 199}
	for _, num := range primes {
		result := isPrime(num)
		if result == false {
			t.Errorf("%d should be a prime number", num)
		}
	}
}

func TestNotPrime(t *testing.T) {
	nonPrimes := []int{1, 4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20,
		21, 22, 24, 25, 26, 27, 28, 30, 32, 33, 34, 35, 36, 38, 39, 40, 42, 44, 45, 46, 48, 49, 50,
		51, 52, 54, 55, 56, 57, 58, 60, 62, 63, 64, 65, 66, 68, 69, 70, 72, 74, 75, 76, 77, 78, 80,
		81, 82, 84, 85, 86, 87, 88, 90, 91, 92, 93, 94, 95, 96, 98, 99, 100}
	for _, num := range nonPrimes {
		result := isPrime(num)
		if result == true {
			t.Errorf("%d should not be a prime number", num)
		}
	}
}
