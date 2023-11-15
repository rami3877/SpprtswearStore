package handerstruct 

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"reflect"
)

func  StructToByte(a any) []byte {
	buffer := bytes.Buffer{}
	if !IsStruct(a) {
		log.Println("not struct")
		return nil
	}
	FieldData := GetDataOfFieldOfStruct(a)
	codeingTobytes := func(a any) {
		binary.Write(&buffer, binary.LittleEndian, a)
	}
	for _, v := range FieldData {
		vR := reflect.ValueOf(v)
		if vR.Kind() == reflect.Slice {
			if vR.Len() == 0 {
				continue
			}
			if index := vR.Index(0); index.Kind() == reflect.String {
				for i := 0; i < vR.Len(); i++ {
					for j := 0; j < index.Len(); j++ {
						codeingTobytes(reflect.ValueOf(v).Index(i).Index(j).Interface())
						// here
					}
				}
			} else {
				for i := 0; i < vR.Len(); i++ {
					codeingTobytes(reflect.ValueOf(v).Index(i).Interface())
				}
			}

		} else if vR.Kind() == reflect.String {
			for _, byteOFstring := range v.(string) {
				codeingTobytes(byteOFstring)
			}
		} else {
			codeingTobytes(v)
		}

	}
	return buffer.Bytes()
}

func  GetSizeOfStruct(a any) (size int, e error) {
	if !IsStruct(a) {
		return 0, errors.New("not struct")
	}
	fieldSize := GetDataOfFieldOfStruct(a)
	if fieldSize == nil {
		return 0, errors.New("cant found size")
	}

	kindOfinedx := func(a any) reflect.Kind {
		return reflect.ValueOf(a).Index(0).Kind()
	}
	ArrayAndSile := func(a reflect.Kind) bool {
		return (a == reflect.Array || a == reflect.Slice)
	}

	for _, v := range fieldSize {
		kind := reflect.ValueOf(v).Kind()

		if kind == reflect.String {
			size += reflect.ValueOf(v).Len()
		} else if ArrayAndSile(kind) && kindOfinedx(v) == reflect.String {
			lenArr := reflect.ValueOf(v).Len()

			for i := 0; lenArr > i; i++ {
				size += reflect.ValueOf(v).Index(i).Len()
			}

		} else if ArrayAndSile(kind) {
			size += int(reflect.ValueOf(v).Index(0).Type().Size()) * int(reflect.ValueOf(v).Len())
		} else {
			size += int(reflect.ValueOf(v).Type().Size())
		}

	}

	return size, e
}

func  StructToRowColm(a any) (err error, data map[string]interface{}) {

	if !IsStruct(a) {
		return errors.New("Not struct"), nil
	}
	data = make(map[string]interface{})
	filedName := GetStructNameOfField(a)
	fieldData := GetDataOfFieldOfStruct(a)
	for i := 0; i < GetNumberFieldOFStruct(a); i++ {
		data[filedName[i]] = fieldData[i]
	}

	return err, data
}

func IsStruct(a any) bool {
	if reflect.TypeOf(a).Kind() != reflect.Struct {
		return false
	}
	return true
}

func  GetStructNameOfField(a any) (field []string) {
	if !IsStruct(a) {
		return nil
	}
	for i := 0; i < GetNumberFieldOFStruct(a); i++ {
		field = append(field, reflect.TypeOf(a).Field(i).Name)
	}
	return field
}

func  GetStructTypeOfFieldString(a any) (field []string) {
	if !IsStruct(a) {
		return nil
	}
	for i := 0; i < GetNumberFieldOFStruct(a); i++ {
		field = append(field, reflect.ValueOf(a).Field(i).Type().String())
	}
	return field
}

func  GetStructTypeOfFieldValue(a any) (field []reflect.Value) {
	if !IsStruct(a) {
		return nil
	}
	for i := 0; i < GetNumberFieldOFStruct(a); i++ {
		field = append(field, reflect.ValueOf(a).Field(i))
	}
	return field
}

func  GetDataOfFieldOfStruct(a any) (data []interface{}) {
	if !IsStruct(a) {
		return nil
	}
	for i := 0; i < GetNumberFieldOFStruct(a); i++ {
		var s interface{} = reflect.ValueOf(a).Field(i).Interface()
		data = append(data, s)
	}
	return data
}
func  GetNumberFieldOFStruct(a any) int {
	return reflect.ValueOf(a).NumField()
}
func GetLenArrayOrSlice(a any) (int, error) {
	if reflect.ValueOf(a).Kind() == reflect.Array || reflect.ValueOf(a).Kind() == reflect.Slice || reflect.String == reflect.ValueOf(a).Kind() {
		return reflect.ValueOf(a).Len(), nil
	}

	return 0, errors.New("Not Slice or Array-> " + reflect.ValueOf(a).Kind().String())
}

func IsArrayOrStringOrSliceStruct(a any) bool {
	return reflect.ValueOf(a).Kind() == reflect.Array || reflect.ValueOf(a).Kind() == reflect.Slice || reflect.ValueOf(a).Kind() == reflect.String
}

func IsArrayOrStringOrSlice(a reflect.Kind) bool {
	return a == reflect.Array || a == reflect.Slice || a == reflect.String
}
// chack if all field of struct 
// is fixed-size
func  IsFxidSize(a any) bool {
	for _, v := range GetStructTypeOfFieldValue(a) {
		if reflect.String == v.Kind() || reflect.Slice == v.Kind() {
			return false
		}
		if reflect.Array == v.Kind() && reflect.String == v.Index(0).Kind() {
			return false
		}

	}

	return true
}
// Get Size Interface 
func SizeOfInterface(a any) (size int) {
	if IsArrayOrStringOrSliceStruct(a) {
		lens, err := GetLenArrayOrSlice(a)
		if err != nil {
			log.Fatal(err)
		}
		if reflect.ValueOf(a).Index(0).Kind() == reflect.String {
			for i := 0; i < lens; i++ {
				size += reflect.ValueOf(a).Index(i).Len()

			}
			return size
		} else if reflect.ValueOf(a).Kind() == reflect.String {
			return lens
		} else {
			size = lens * int(reflect.ValueOf(a).Index(0).Type().Size())
			return size
		}

	}
	return int(reflect.ValueOf(a).Type().Size())
}

/*
	switch a.(type) {
	case *bool:

	case *int8:

	case *uint8:

	case *int16:

	case *uint16:

	case *int32:

	case *uint32:

	case *int64:

	case *uint64:

	case *float32:

	case *float64:

	case []bool:
		for i, x := range bs { // Easier to loop over the input for 8-bit values.

		}
	case []int8:
		for i, x := range bs {

		}
	case []uint8:
	case []int16:
		for i := range

		}
	case []uint16:
		for i := range

		}
	case []int32:
		for i := range

		}
	case []uint32:
		for i := range

		}
	case []int64:
		for i := range

		}
	case []uint64:
		for i := range

		}
	case []float32:
		for i := range

		}
	case []float64:
		for i := range

		}
	default:
	}


*/
