package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/votes/config"
	"github.com/votes/model"
	"log"
	"net/http"
	jwt "github.com/appleboy/gin-jwt"
)


type User = model.User

func GetUser(c *gin.Context) {
	var user User
	uuid := c.Param("uuid")
	config.DB.Where("uuid = ?", uuid).Find(&user)
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var u User
	err := c.BindJSON(&u)
	claims := jwt.ExtractClaims(c)
	var accessLevel int = int(claims["accessLevel"].(float64))
	if accessLevel == 0 && u.AccessLevel == 1 {
		c.JSON(http.StatusUnauthorized, errors.New("You need to be an admin to create an admin"))
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
	}
	if u.IsValid() != nil {
		log.Println(u.IsValid())
		c.JSON(http.StatusBadRequest, u.IsValid())
	}
	config.DB.NewRecord(u)
	config.DB.Create(&u)
	c.JSON(http.StatusOK, u)
}

func UpdateUser(c *gin.Context) {
	var u User
	uuid := c.Param("uuid")
	config.DB.Where("uuid = ?", uuid).Find(&u)
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
	fmt.Println(updatedUser)
	u.SetFirstname(updatedUser.FirstName)
	u.SetLastname(updatedUser.LastName)
	u.SetEmail(updatedUser.Email)
	u.SetPassword(updatedUser.Password)
	u.DateOfBirth = updatedUser.DateOfBirth
	config.DB.Save(&u)
	c.JSON(http.StatusOK, u)
}

func DeleteUser(c *gin.Context) {
	var u User
	uuid := c.Param("uuid")
	config.DB.Where("uuid = ?", uuid).Find(&u)
	config.DB.Delete(&u)
	c.JSON(http.StatusOK, u)
}
