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
	timestart := time.Now() // used to keep track of the implementation processing time

	input, err := os.ReadFile("input.json") // reading the .json file

	if err != nil {
		panic(err)
	} // checking for errors

	var parsed any

	// we pass input to Unmarshal, which changes the second value passed, which must be a pointer, otherwise, it will not be able to change the actual value.
	if err := json.Unmarshal(input, &parsed); err != nil {
		panic(err)
	}
	output, err := transformMap(parsed)

	// checking for errors
	if err != nil {
		panic(err)
	}

	// indenting the printout with MarshalIndent to make it readable
	printOut, err := json.MarshalIndent([]any{output}, "", "  ")

	// checking for errors
	if err != nil {
		panic(err)
	}

	// printing out the transformed result
	fmt.Println(string(printOut))

	// reporting the processing time
	timeFinish := time.Since(timestart)

	fmt.Println(timeFinish)

}

// transforms and validates string values
func transformString(input any) (any, error) {

	// type assertion to check if the value is a string
	str, ok := input.(string)

	if !ok {
		return nil, errors.New("invalid string type")
	}

	// trimming whitespace from the string
	str = strings.TrimSpace(str)

	// checking whether it is an empty string
	if str == "" {
		return nil, errors.New("empty string")
	}

	// type assertion to check whether the string is a RFC3339 formatted string
	t, err := time.Parse(time.RFC3339, str)

	// returning the time formatted string in the desired data type
	if err == nil {
		return t.Unix(), nil
	}

	// returning transformed string after validation
	return str, nil

}

// transforms and validates numeric values
func transformNum(input any) (any, error) {

	// type assertion to check if the value is a string
	num, ok := input.(string)

	if !ok {
		return nil, errors.New("invalid input")
	}

	// trimming whitespace from the string
	num = strings.TrimSpace(num)

	// checking whether the value contains an integer data type
	intVar, err := strconv.Atoi(num)

	// early return in case the value is an integer
	if err == nil {
		return intVar, nil
	}

	// checking whether the value contains a float data type
	floatVar, err := strconv.ParseFloat(num, 64)

	// if the string contains neither integer nor float, return an error
	if err != nil {
		return nil, errors.New("invalid number")
	}

	// otherwise, return the float data type
	return floatVar, nil
}

// transforms and validates boolean values
func transformBool(input any) (bool, error) {

	// type assertion to check if the value is a string
	boolVar, ok := input.(string)

	// if the value is not a string, return an error
	if !ok {
		return false, errors.New("invalid input")
	}

	// trimming whitespace from the string
	boolVar = strings.TrimSpace(boolVar)

	// checking whether the value contains a valid boolean type
	boolResult, err := strconv.ParseBool(boolVar)

	if err == nil {
		return boolResult, nil
	}

	// if the data type is invalid, it will return an error
	return false, err
}

// transforms and validates null values
func transformNull(input any) (any, error) {

	// type assertion to check if the value is a string
	nullVar, ok := input.(string)

	if !ok {
		return nil, errors.New("invalid input")
	}

	// trimming whitespace from the string
	nullVar = strings.TrimSpace(nullVar)

	// checking whether the value contains a valid boolean expression
	nullResult, err := strconv.ParseBool(nullVar)

	// if the boolean result is true, return nil (null) without an error
	if err == nil && nullResult == true {
		return nil, nil
	}

	// after validation, return an error because the value is not null
	return nil, errors.New("not null")
}

// transforms and validates list types
func transformList(input any) ([]any, error) {

	// type assertion to check if the value is a slice
	l, ok := input.([]any)

	if !ok {
		return nil, errors.New("not a list")
	}

	// creating new slice to store the transformed values
	output := make([]any, 0)

	// ranging over each index to transform each value
	for _, value := range l {

		// indirect recursive call to transform to process the keys
		transformed, err := transform(value)

		// break out of for loop in case of error
		if err != nil {
			continue
		}

		// adding the transformed value to the output slice
		output = append(output, transformed)
	}

	// if the list is empty, it will be omitted
	if len(output) == 0 {
		return nil, errors.New("empty list")
	}

	// after validation and transformation, return the output slice
	return output, nil
}

// transforms and validates map types
func transformMap(input any) (map[string]any, error) {

	// type assertion to check if the value is a map
	m, ok := input.(map[string]any)

	if !ok {
		return nil, errors.New("not a map")
	}

	// creating a new map to store the transformed values
	output := make(map[string]any)

	// range over each key, sanitize it, and transform each value
	for innerKey, value := range m {
		trimmedKey := strings.TrimSpace(innerKey)

		if trimmedKey == "" {
			continue
		}

		// indirect recursive call to the transform function
		transformed, err := transform(value)
		if err != nil {
			continue
		}

		// assign newly transformed values to the output map
		output[trimmedKey] = transformed
	}

	// checking for an empty map
	if len(output) == 0 {
		return nil, errors.New("empty map")
	}

	// returning the newly created map with the transformed values
	return output, nil
}

// transform checks each "type key" "by using type checking in a switch statement
func transform(input any) (any, error) {
	switch n := input.(type) {
	case map[string]any:
		for key, value := range n {

			// sanitizing the keys of the leading and trailing whitespace
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

	// mandatory return statement if all else fails, that includes an error
	return nil, errors.New("no valid data type")
}
