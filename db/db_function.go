package db

import (
	"db/file"
	handerstruct "db/handerStruct"
	"errors"
	"log"
	"os"
	"strings"
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

func (db *DB) InsterStruct(a any, name string) error {
	if !handerstruct.IsStruct(a) {
		return errors.New("is not struct")
	} else if len(name) == 0 {
		return errors.New("no name ")
	}
	db.FileDataBase.Seek(0, 0)
	Node := file.NewFileNodeFormFile(db.FileDataBase)
	root := Node.HeadFile
	var tableNode *file.Node = nil
	for root.Next != nil {
		if root.Name == name {
			tableNode = root
			break
		}
		root = root.Next
	}
	if tableNode == nil {
		return errors.New("ERROR table not exit {" + name + "} try use CreateTable function ")
	}

	mapKT := file.UncodeFromDataBase(tableNode.Data)
	for k := range mapKT {
		if !handerstruct.FoundField(a, k) {
			return errors.New("ERROR: cant found " + k)
		}
		// TODO  check if type of field  some field type of pass struct 
	}

	return nil
}

func (db *DB) CreateTable(a any, name string) error {
	name = strings.TrimSpace(name)
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
		db.FileDataBase.Seek(0, 0)
		nodeFile := file.InitFileNode()
		nodeFile.InsrtByIndex(0, name, []byte(handerstruct.JoinNameOFFieldAndType(a)))
		nodeFile.WritToFile(db.FileDataBase)
	} else {
		db.FileDataBase.Seek(0, 0)
		nodeFile := file.NewFileNodeFormFile(db.FileDataBase)
		root := nodeFile.HeadFile
		for i := 0; i != nodeFile.Len; i++ {
			if root.Name == name {
				return errors.New("is exitd " + name)
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

func (db *DB) Close() {
	db.FileDataBase.Close()
}
