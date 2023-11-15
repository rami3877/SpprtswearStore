package main

import (
	"db"
	"log"
)

func main() {
	dataBase, err := db.OpenDB("rami")
	if err != nil {
		log.Fatal(err)
	}
	type s struct {
		Name [40]byte
		B    uint64
	}
	_ = dataBase
	if err := dataBase.CreateTable(s{}, "user"); err != nil {
		log.Println(err)
	}

	if err := dataBase.CreateTable(s{}, "user1"); err != nil {
		log.Println(err)
	}

}
