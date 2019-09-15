package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/votes/config"
	"github.com/votes/model"
	"log"
	"net/http"
)


type User = model.User

func GetUser(c *gin.Context) {
	var user User
	uuidParam := c.Query("uuid")
	config.DB.Where("uuid = ?", uuidParam).Find(&user)
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var u User
	err := c.BindJSON(&u)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	u.SetPassword(u.Password)
	config.DB.NewRecord(u)
	config.DB.Create(&u)
	fmt.Println(u)
	c.JSON(http.StatusOK, u)
}
