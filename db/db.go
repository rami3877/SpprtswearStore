package db

import (
	"errors"
	"log"
	"os"
	"slices"
	"structs"
	"go.etcd.io/bbolt"
)

type DataBase struct {
	NameDataBase string
	InStock      stock
	OutStock     stock
	Orders       order
	Users        user
}


var MainDB *DataBase 

var (
	ErrDataBaseOutStock = errors.New("Sorry, out of stock")
)

func (db *DataBase) AddCommint(id int, Container, kind string, commint structs.UserCommint) error {
	err := db.Users.GetUser(commint.Username, nil)
	if err == ErrUsereNotFound {
		return err
	}
	return db.InStock.addCommint(id, Container, kind, commint)
}

func (db *DataBase) Buy(order Orders) error {
	kinds, err := db.InStock.GetModelsInKind(order.IdModel, 0, order.Container, order.Kind)
	if err != nil {
		return err
	}
	err = db.Users.GetUser(order.Username, nil)
	if err == ErrUsereNotFound {
		return err
	}

	if len(kinds) == 0 {
		log.Fatal("[error in code ] Buy database len is  zero")
	}

	kind := kinds[0]
	index := -1
	for k, color := range kind.Sizes[order.Size].Colors {
		if color.ColorName == order.Color {
			index = k
		}
	}
	if index == -1 {
		return ErrDataBaseOutStock
	}
	kind.Sizes[order.Size].Colors[index].Qty--
	if kind.Sizes[order.Size].Colors[index].Qty <= 0 {
		// i dont know if will work
		news := kind.Sizes[order.Size].Colors
		news = slices.DeleteFunc(news, func(s structs.Color) bool {
			return s.ColorName == order.Color
		})
		db.InStock.UpdataSizeFromModel(order.IdModel, order.Container, order.Kind, order.Size, kind.Sizes[order.Size])

	} else {
		db.InStock.UpdataSizeFromModel(order.IdModel, order.Container, order.Kind, order.Size, kind.Sizes[order.Size])

	}

	if len(kind.Sizes) == 0 {
		_, err := db.OutStock.ContainerAndKindIsExited(order.Container, order.Kind)
		if err != nil {
			db.OutStock.AddNewContainer(order.Container)
			db.OutStock.NewKindtoContainer(order.Container, order.Kind)
			db.OutStock.AddModelToKind(order.Container, order.Kind, &kind, false)
		}
		db.InStock.DeleteModelFromKind(order.Container, order.Kind, order.IdModel)

	} else if len(kind.Sizes[order.Size].Colors) == 0 {
		db.InStock.DeleteSizeFromModel(order.IdModel, order.Container, order.Kind, order.Size)
	}
	db.Orders.Add(order)

	return nil
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
	MainDB = &db
	return &db
}


