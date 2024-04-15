package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	timestart := time.Now()

	input, err := os.ReadFile("input.json")

	if err != nil {
		panic(err)
	}

	var parsed any

	//we pass input to unmarshal, which changes the second value passed, which must be a pointer, otherwise, it will not persist.
	if err := json.Unmarshal(input, &parsed); err != nil {
		panic(err)
	}
	output, err := transformMap(parsed)

	if err != nil {
		panic(err)
	}
	printOut, err := json.MarshalIndent([]any{output}, "", "  ")

	if err != nil {
		panic(err)
	}

	fmt.Println(string(printOut))

	timeFinish := time.Since(timestart)

	fmt.Println(timeFinish)

}

func transformString(input any) (any, error) {
	str, ok := input.(string)

	if !ok {
		return nil, errors.New("invalid string type")
	}

	str = strings.TrimSpace(str)

	if str == "" {
		return nil, errors.New("empty string")
	}

	t, err := time.Parse(time.RFC3339, str)

	if err == nil {
		return t.Unix(), nil
	}

	return str, nil

}

func transformNum(input any) (any, error) {
	num, ok := input.(string)

	if !ok {
		return nil, errors.New("invalid input")
	}

	num = strings.TrimSpace(num)

	intVar, err := strconv.Atoi(num)

	if err == nil {
		return intVar, nil
	}

	floatVar, err := strconv.ParseFloat(num, 64)

	if err != nil {
		return nil, errors.New("invalid number")
	}
	return floatVar, nil
}

func transformBool(input any) (bool, error) {
	boolVar, ok := input.(string)

	if !ok {
		return false, errors.New("invalid input")
	}

	boolVar = strings.TrimSpace(boolVar)

	boolResult, err := strconv.ParseBool(boolVar)

	if err == nil {
		return boolResult, nil
	}

	return false, err
}

func transformNull(input any) (any, error) {
	nullVar, ok := input.(string)

	if !ok {
		return nil, errors.New("invalid input")
	}

	nullVar = strings.TrimSpace(nullVar)

	nullResult, err := strconv.ParseBool(nullVar)

	if err == nil && nullResult == true {
		return nil, nil
	}

	return nil, errors.New("not null")
}

func transformList(input any) ([]any, error) {
	l, ok := input.([]any)

	if !ok {
		return nil, errors.New("not a list")
	}

	output := make([]any, 0)

	for _, value := range l {

		transformed, err := transform(value)
		if err != nil {
			continue
		}

		output = append(output, transformed)
	}

	if len(output) == 0 {
		return nil, errors.New("empty list")
	}

	return output, nil
}

func transformMap(input any) (map[string]any, error) {
	m, ok := input.(map[string]any)

	if !ok {
		return nil, errors.New("not a map")
	}

	output := make(map[string]any)

	for innerKey, value := range m {
		trimmedKey := strings.TrimSpace(innerKey)

		if trimmedKey == "" {
			continue
		}
		transformed, err := transform(value)
		if err != nil {
			continue
		}
		output[trimmedKey] = transformed
	}

	if len(output) == 0 {
		return nil, errors.New("empty map")
	}

	return output, nil
}

func transform(input any) (any, error) {
	switch n := input.(type) {
	case map[string]any:
		for key, value := range n {
			key = strings.TrimSpace(key)
			switch key {
			case "S":
				return transformString(value)
			case "N":

				return transformNum(value)
			case "BOOL":
				return transformBool(value)
			case "NULL":
				return transformNull(value)
			case "L":
				return transformList(value)
			case "M":
				return transformMap(value)
			}
		}
	}
	return nil, errors.New("no valid data type")
}
