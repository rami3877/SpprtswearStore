package handerstruct

import "reflect"

func GetDataOfFieldOfStructByName(a any, name string) (data interface{}) {
	if !IsStruct(a) {
		return nil
	}
	return reflect.ValueOf(a).FieldByName(name).Interface()
}
func FoundField(a any, nameField string) bool {

	if !IsStruct(a) {
		return false
	}
	if reflect.Invalid == reflect.ValueOf(a).FieldByName(nameField).Kind() {
		return false
	}

	return true
}

func SomeType(a any, fieldName string, typeField string) bool {
	if reflect.ValueOf(a).FieldByName(fieldName).Type().String() != typeField {
		return false
	} else {
		return true
	}
}
