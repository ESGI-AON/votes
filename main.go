package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/votes/config"
	"github.com/votes/controller"
	"github.com/votes/model"
)

type User = model.User
type Vote = model.Vote

var err error

func main(){

	config.DB, err = gorm.Open("postgres", "host=localhost port=5432 user=root password=root dbname=govotes sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}

	config.DB.AutoMigrate(&User{})

	r := gin.Default()
	r.GET("/user", controller.GetUser)
	r.POST("/user", controller.CreateUser)
	r.PUT("/user", controller.UpdateUser)
	r.DELETE("/user", controller.DeleteUser)

	config.DB.AutoMigrate(&Vote{})

	r.GET("/vote", controller.GetVote)
	r.POST("/vote", controller.CreateVote)
	r.PUT("/vote", controller.UpdateVote)
	r.DELETE("/vote", controller.DeleteVote)

	r.Run(":8080")

	defer config.DB.Close()
}