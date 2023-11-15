package handerstruct

import "reflect"

func FoundField(a any, nameField string) bool {
	if !IsStruct(a) {
		return false
	}
	if reflect.Invalid == reflect.ValueOf(a).FieldByName(nameField).Kind() {
		return false
	}

	return true
}
