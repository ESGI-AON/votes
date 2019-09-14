package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/votes/model"
	"log"
	"time"
)

type User = model.User

func main(){
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=root password=root dbname=govotes sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(time.Date(1994, 9, 22, 0, 0, 0, 0, time.UTC))
	return
	h := sha256.New()
	h.Write([]byte("toto"))
	hashedPwd := base64.URLEncoding.EncodeToString(h.Sum(nil))

	user := &User{
		AccessLevel: 0,
		FirstName:   "Alex",
		LastName:    "Tea",
		Email:       "alextea2@gmail.com",
		Password:    hashedPwd,
		DateOfBirth: time.Now(),
	}
	db.AutoMigrate(&User{})

	if len(user.IsValid()) > 0 {
		log.Println(user.IsValid())
	} else {
		db.Create(user)
	}

	fmt.Println(err)
	defer db.Close()
}