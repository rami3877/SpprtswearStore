package main

import (
	"fmt"
	"reflect"
)

func main() {
	type d struct {
		Name int64
	}

	if reflect.Invalid == reflect.ValueOf(d{}).FieldByName("N").Kind() {
		fmt.Println("no")
	}
}
