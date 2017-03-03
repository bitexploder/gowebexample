package model

import (
	"github.com/asdine/storm"
	"time"
)

type User struct {
	ID       int `storm:"id,increment"`
	Created  time.Time
	Name     string
	Email    string
	Username string
	PassHash string
	IsAdmin  bool
}

func GetUser(db *storm.DB, userId int) (User, error) {
	var user User
	err := db.One("ID", userId, &user)
	return user, err
}

func GetUserByUsername(db *storm.DB, username string) (User, error) {
	var user User
	err := db.One("Username", username, &user)
	return user, err
}

func GetUsers(db *storm.DB) ([]User, error) {
	var users []User
	err := db.All(&users)
	return users, err
}

func CreateUser(db *storm.DB, user User) error {
	err := db.Save(&user)
	return err
}

func UpdateUser(db *storm.DB, user User) error {
	err := db.Save(&user)
	return err
}

func DeleteUser(db *storm.DB, userId int) error {
	u, err := GetUser(db, userId)
	if err != nil {
		return err
	}
	err = db.DeleteStruct(&u)
	return err
}
