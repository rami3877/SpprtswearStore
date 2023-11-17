package file

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

func UncodeFromDataBase(data []byte) map[string]string {
	types := make(map[string]string)

	key := ""
	Type := ""

	for i := 0; i < len(data); i++ {
		if data[i] == ':' {
			i++
			for ; i < len(data); i++ {
				if data[i] == '\r' {
					types[key] = Type
					key = ""
					Type = ""
					break
				}
				Type += string(data[i])
			}
		}
		if data[i] != '\r' {
			key += string(data[i])
		}

	}
	return types
}

// dont use map
func CodeingToRow(a any) []byte {
	buffer := bytes.NewBuffer(nil)
	if reflect.ValueOf(a).Kind() == reflect.Map {
		return nil
	}
	if reflect.ValueOf(a).Kind() == reflect.String {

		stringVaule := a.(string)
		buffer.WriteString(stringVaule)

	} else if reflect.ValueOf(a).Kind() == reflect.Array || reflect.ValueOf(a).Kind() == reflect.Slice {
		binary.Write(buffer, binary.LittleEndian, uint64(reflect.ValueOf(a).Len()))
		if reflect.ValueOf(a).Index(0).Kind() == reflect.String {
			for i := 0; i < reflect.ValueOf(a).Len(); i++ {
				buffer.WriteString(reflect.ValueOf(a).Index(i).Interface().(string))
				buffer.WriteByte(0)
			}
		} else {
			 // bad storage  when fixed array be byte 
			for i := 0; i < reflect.ValueOf(a).Len(); i++ {
				if reflect.ValueOf(a).Index(0).Kind() == reflect.Int {
					binary.Write(buffer, binary.LittleEndian, int64(reflect.ValueOf(a).Index(i).Interface().(int)))
				} else if reflect.ValueOf(a).Index(0).Kind() == reflect.Uint {
					binary.Write(buffer, binary.LittleEndian, uint64(reflect.ValueOf(a).Index(i).Interface().(uint)))
				} else {
					binary.Write(buffer, binary.LittleEndian, reflect.ValueOf(a).Index(i).Interface())
				}
			}

		}
	} else {
		if reflect.ValueOf(a).Kind() == reflect.Int {
			binary.Write(buffer, binary.LittleEndian, int64(a.(int)))
		} else if reflect.ValueOf(a).Kind() == reflect.Uint {
			binary.Write(buffer, binary.LittleEndian, uint64(reflect.ValueOf(a).Interface().(uint)))
		} else {
			binary.Write(buffer, binary.LittleEndian, a)
		}
	}

	return buffer.Bytes()
}
