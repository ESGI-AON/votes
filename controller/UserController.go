package controller

import (
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
	if u.IsValid() != nil {
		log.Println(u.IsValid())
		c.JSON(http.StatusBadRequest, u.IsValid())
		return
	}
	config.DB.NewRecord(u)
	config.DB.Create(&u)
	c.JSON(http.StatusOK, u)
}

func UpdateUser(c *gin.Context) {
	var u User
	uuidParam := c.Query("uuid")
	config.DB.Where("uuid = ?", uuidParam).Find(&u)
	var updatedUser User
	err := c.BindJSON(&updatedUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if len(u.IsValid()) > 0 {
		log.Println(u.IsValid())
		c.JSON(http.StatusBadRequest, u.IsValid())
		return
	}
	u.FirstName = updatedUser.FirstName
	u.LastName = updatedUser.LastName
	u.Email = updatedUser.Email
	u.SetPassword(updatedUser.Password)
	u.DateOfBirth = updatedUser.DateOfBirth
	config.DB.Save(&u)
	c.JSON(http.StatusOK, u)
}