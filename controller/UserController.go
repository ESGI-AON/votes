package controller

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/votes/config"
	"github.com/votes/helper"
	"github.com/votes/model"
	"log"
	"net/http"
)


type User = model.User

// get user by uuid
func GetUser(c *gin.Context) {
	var user User
	uuid := c.Param("uuid")
	config.DB.Where("uuid = ?", uuid).Find(&user)
	c.JSON(http.StatusOK, user)
}

// create new user
func CreateUser(c *gin.Context) {
	var u User
	err := c.BindJSON(&u)
	if u.FirstName == "" || u.LastName == "" || u.Password == "" {
		c.JSON(http.StatusBadRequest, "Firstname, Lastname, Password are required")
		return
	}
	claims := jwt.ExtractClaims(c)
	accessLevel := helper.GetAccessLevel(claims)
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
		strErrors := make([]string, len(u.IsValid()))
		for i, err := range u.IsValid() {
			strErrors[i] = err.Error()
		}
		c.JSON(http.StatusBadRequest, strErrors)
		return
	}
	config.DB.NewRecord(u)
	config.DB.Create(&u)
	c.JSON(http.StatusOK, u)
}


// update existing user by uuid
func UpdateUser(c *gin.Context) {
	var u User
	uuid := c.Param("uuid")
	config.DB.Where("uuid = ?", uuid).Find(&u)
	var updatedUser User
	err := c.BindJSON(&updatedUser)
	claims := jwt.ExtractClaims(c)
	accessLevel := helper.GetAccessLevel(claims)
	jwtUUID :=  helper.GetUUID(claims)
	if (u.UUID != jwtUUID) && accessLevel == 0 {
		c.JSON(http.StatusUnauthorized, "You need to an admin to edit another user")
		return
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if len(updatedUser.IsValid()) > 0 {
		strErrors := make([]string, len(updatedUser.IsValid()))
		for i, err := range updatedUser.IsValid() {
			strErrors[i] = err.Error()
		}
		c.JSON(http.StatusBadRequest, strErrors)
		return
	}
	fmt.Println(updatedUser)
	u.SetFirstname(updatedUser.FirstName)
	u.SetLastname(updatedUser.LastName)
	u.SetEmail(updatedUser.Email)
	u.SetPassword(updatedUser.Password)
	u.SetBirthDate(updatedUser.DateOfBirth)
	config.DB.Save(&u)
	c.JSON(http.StatusOK, u)
}


// soft delete user by uuid
func DeleteUser(c *gin.Context) {
	var u User
	uuid := c.Param("uuid")
	claims := jwt.ExtractClaims(c)
	accessLevel := helper.GetAccessLevel(claims)
	if accessLevel == 0 {
		c.JSON(http.StatusUnauthorized, "You need to be an admin to delete a user")
		return
	}
	config.DB.Where("uuid = ?", uuid).Find(&u)
	config.DB.Delete(&u)
	c.JSON(http.StatusOK, u)
}
