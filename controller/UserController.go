package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	googleUUID "github.com/google/uuid"
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
		c.JSON(http.StatusUnauthorized, "You need to be an admin to create an admin")
		return
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
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
	uuid := c.Param("uuid")
	config.DB.Where("uuid = ?", uuid).Find(&u)
	var updatedUser User
	err := c.BindJSON(&updatedUser)
	claims := jwt.ExtractClaims(c)
	var accessLevel int = int(claims["accessLevel"].(float64))
	jwtUUID, err :=  googleUUID.Parse(fmt.Sprintf("%v", claims["uuid"]))
	if (u.UUID != jwtUUID) && accessLevel == 0 {
		c.JSON(http.StatusUnauthorized, "You need to an admin to edit another user")
		return
	}
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
