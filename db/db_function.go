package db

import (
	handerstruct "db/HanderStruct"
	"db/tools"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"unsafe"
)

func OpenDB(name string) (*DB, error) {
	file, err := os.OpenFile(name+".db", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	return &DB{
		FileDataBase: file,
	}, err
}

func (db *DB) CreateTable(a any, name string) error {

	if len(name) > 50 {
		return errors.New("name of table so long")
	}
	if !handerstruct.IsStruct(a) {
		return errors.New("not struct")
	}

	FileDB := "DBFILE_" + db.FileDataBase.Name()
	os.Mkdir(FileDB, 0777)
	info, _ := db.FileDataBase.Stat()

	file, err := os.OpenFile(FileDB+"/"+name+".db", os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if info.Size() == 0 {
		if !handerstruct.IsFxidSize(a) {
			return errors.New("not fixe")
		}
		buffer :=tools.NewBuffer(nil)
		// how many table -> headerSize
		buffer.Ints(int64(1))
		var nameByte [50]byte
		copy(nameByte[:], []byte(name))
		buffer.Write(nameByte[:])
		str := db.JoinNameOFFieldAndType(a)
		buffer.Ints(int64(len(str)))
		buffer.WriteString(str)
		buffer.WriteToFile(db.FileDataBase)
	} else {
		if !handerstruct.IsFxidSize(a) {
			return errors.New("not fixe")
		}
		db.seekTo(0)
		buffer := tools.NewBuffer(nil)
		buffer.ReadFile(db.FileDataBase)
		lenHeader := binary.LittleEndian.Uint64(buffer.Next(8))
		for i := 0; i < int(lenHeader); i++ {
			str := buffer.Next(50)
			s := 0

			for j := 0; j < len(name); j++ {
				if str[j] == []byte(name)[j] {
					s++
				}
				if j+1 == len(name) && j+1 < len(name) && str[j+1] != 0 {
					s = 0
					break
				}
			}

			if s == len(name) {
				return errors.New("is table Exit")
			}

			buffer.Next(8)
		}
		buffer.WriteFromIndex(binary.LittleEndian.AppendUint64(nil, lenHeader+1), 0)
		var nameByte [50]byte
		copy(nameByte[:], []byte(name))
		buffer.AppendFromIdex((int(lenHeader) * 58), nameByte[:])

		str := db.JoinNameOFFieldAndType(a)

		buffer.AppendFromIdex((int(lenHeader*(58)) + 50-16), tools.Ints( int64 (len(str))))
		buffer.WriteString(str)
		g := db.FileDataBase.Name()
		os.Remove(db.FileDataBase.Name())
		db.FileDataBase , _ = os.OpenFile(g, os.O_CREATE| os.O_RDWR, 0600)
		buffer.WriteToFile(db.FileDataBase)


	}

	return nil
}

func (db *DB) GetDataAsStruct(a any) {

}

func (db *DB) ReadIndex() (it []indexTable) {
	db.seekTo(0)
	var d int64
	db.ReadData(&d)
	for i := int64(0); i != d; i++ {
		indexes := new(indexTable)
		db.ReadData(indexes)
		it = append(it, *indexes)
	}
	return it
}

func (db *DB) ReadIndex2() (it []indexTable) {
	db.seekTo(0)
	var d int64
	for {
		if err := binary.Read(db.FileDataBase, binary.LittleEndian, &d); err != nil {
			break
		}
		fmt.Println(IntToByteArray(d))
	}
	return it
}

func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}
func (db *DB) ReadData(a any) {
	//db.seek = int64(SizeOfInterface(a))
	binary.Read(db.FileDataBase, binary.LittleEndian, a)
}

func (db *DB) seekTo(s int64) {
	db.seek, db.errSeek = db.FileDataBase.Seek(s, 0)
	if db.errSeek != nil {
		log.Fatal(db.errSeek)
	}
}

func (db *DB) JoinNameOFFieldAndType(a any) string {
	FieldName := handerstruct.GetStructNameOfField(a)
	TypeName := handerstruct.GetStructTypeOfFieldString(a)
	str := ""
	for i, v := range FieldName {
		str += v + ":" + TypeName[i]
		if i != len(TypeName)-1 {
			str += "\r"
		}
	}
	str += "\r"
	return str
}

func (db *DB) Close() {
	db.FileDataBase.Close()

}
func (db *DB) fetchDataDB(s any) (error, string) {
	return nil, "data"
}

func (db *DB) WriteData(s any) error {

	return nil
}
