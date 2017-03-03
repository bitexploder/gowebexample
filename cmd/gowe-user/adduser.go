package main

import (
	"flag"
	"fmt"

	"github.com/asdine/storm"
	"github.com/bitexploder/gowebexample/model"
	"golang.org/x/crypto/bcrypt"
)

func main() {

	dbPath := flag.String("dbpath", "gowe.db", "database path")
	name := flag.String("name", "First Last", "firstname lastname")
	password := flag.String("password", "", "a password...")
	username := flag.String("username", "", "a username...")
	flag.Parse()

	db, err := storm.Open(*dbPath)
	if *username == "" || *password == "" {
		fmt.Println("Need username and password")
		return
	}

	u := model.User{}

	u, err = model.GetUserByUsername(db, *username)
	if err == nil {
		fmt.Println("Updating user and promoting to admin")
	} else {
		u.Username = *username
		u.Name = *name
		u.IsAdmin = true

		fmt.Println("Creating user and promoting to admin")
	}
	u.IsAdmin = true

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	u.PassHash = string(hash)
	fmt.Printf("%+v\n", u)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	err = model.UpdateUser(db, u)
	if err != nil {
		fmt.Printf("err: %s\n")
		return
	}
}
