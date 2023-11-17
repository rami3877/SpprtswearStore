package main

import (
	"db"
	"fmt"
	"log"
)

func main() {
	dataBase, err := db.OpenDB("rami")
	defer dataBase.Close()
	if err != nil {
		log.Fatal(err)
	}
	type s struct {
		Name [40]byte
		B    uint64
	}
	_ = dataBase

	fg := s{}
	copy(fg.Name[:] , []byte("dasdasdsada"))
	fg.B = 12 
	if err = dataBase.InsterStruct(fg, "user"); err != nil {
		fmt.Println(err)
	}

}
