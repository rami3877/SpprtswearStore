package main

import (
	"db"
	"fmt"
	"log"
	"structs"
)

func main() {
	d := db.OpenDirDataBase("DataBase")
	user := structs.User{}
	user.Password = "dasdasdasdasd"

	err := d.Users.UpdataPassword("ramia", "adminSddasdsad", "thisnewPassword")
	if err != nil {
		 log.Println("da")
		log.Fatal(err)
	}
	userReturn := structs.User{}
	err = d.Users.GetUser("ramia", &userReturn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(userReturn)

	d.Close()

}
