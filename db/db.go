package db

import (
	"go.etcd.io/bbolt"
	"log"
	"os"
)

type DataBase struct {
	NameDataBase string
	InStock      stock
	OutStock     stock
	Orders       order
	Users        user
}

func (db *DataBase) Close() {
	for _, v1 := range db.InStock.database {
		for _, data := range v1 {
			data.Close()
		}
	}
	for _, v1 := range db.OutStock.database {
		for _, data := range v1 {
			data.Close()
		}
	}
	db.Orders.dataBase.Close()
	db.Users.dataBase.Close()
}

func OpenDirDataBase(name string) *DataBase {
	db := DataBase{}
	db.OutStock.pathStock = name + "/outStock"
	db.InStock.pathStock = name + "/inStock"
	err := os.Mkdir(name, 0770)
	if err != nil && !os.IsExist(err) {
		log.Fatal("[database] " + err.Error())
	}
	err = os.Mkdir(db.InStock.pathStock, 0770)
	if err != nil && !os.IsExist(err) {
		log.Fatal("[database] " + err.Error())
	}
	err = os.Mkdir(db.OutStock.pathStock, 0770)
	if err != nil && !os.IsExist(err) {
		log.Fatal("[database] " + err.Error())
	}

	db.InStock.database = make(map[string]map[string]*bbolt.DB)
	db.OutStock.database = make(map[string]map[string]*bbolt.DB)
	db.NameDataBase = name

	if err := db.InStock.init(); err != nil {
		log.Fatal(err)
	}
	if err := db.OutStock.init(); err != nil {
		log.Fatal(err)
	}
	db.Users.dataBase, err = bbolt.Open(name+"/user"+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Users.dataBase.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	// init Orders
	db.Orders.dataBase, err = bbolt.Open(name+"/orders"+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Orders.dataBase.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("orders"))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
	return &db
}
