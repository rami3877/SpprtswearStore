package main

import (
	"fmt"
	"reflect"
)

func main() {
	type d struct {
		Name [12]int64
	}

	if reflect.Invalid == reflect.ValueOf(d{}).FieldByName("N").Kind() {
		fmt.Println("no")
	}


	fmt.Println(reflect.ValueOf(d{}).Field(0).Type().String())

	fmt.Println(reflect.ValueOf([50]byte{}).Index(0).Kind())
}
