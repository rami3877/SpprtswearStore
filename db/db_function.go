package db

import (
	"db/file"
	"db/handerStruct"
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

	f, err := os.OpenFile(FileDB+"/"+name+".db", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	info, err := db.FileDataBase.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if info.Size() == 0 {
		db.seekTo(0)
		nodeFile := file.InitFileNode()
		nodeFile.InsrtByIndex(0, name, []byte(handerstruct.JoinNameOFFieldAndType(a)))
		nodeFile.WritToFile(db.FileDataBase)
	} else {
		db.seekTo(0)
		nodeFile := file.NewFileNodeFormFile(db.FileDataBase)
		root := nodeFile.HeadFile
		for i := 0; i != nodeFile.Len; i++ {
			if root.Name == name {
				return errors.New("is exitd "+ name)
			}
			if root.Next == nil {
				break
			}
			root = root.Next
		}
		db.FileDataBase.Seek(0, 0)
		dataofTable := []byte(handerstruct.JoinNameOFFieldAndType(a))
		nodeFile.Append(&file.Node{Name: name, Data: dataofTable, Size: uint64(len(dataofTable))})
		nodeFile.WritToFile(db.FileDataBase)
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

func (db *DB) Close() {
	db.FileDataBase.Close()

}
func (db *DB) fetchDataDB(s any) (error, string) {
	return nil, "data"
}

func (db *DB) WriteData(s any) error {

	return nil
}
