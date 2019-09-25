package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/votes/config"
	"github.com/votes/model"
	"log"
	"net/http"
)


type Vote = model.Vote

func GetVote(c *gin.Context) {
	var vote Vote
	uuidParam := c.Query("uuid")
	config.DB.Where("uuid = ?", uuidParam).Find(&vote)
	c.JSON(http.StatusOK, vote)
}

func CreateVote(c *gin.Context) {
	var v Vote
	err := c.BindJSON(&v)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	config.DB.NewRecord(v)
	config.DB.Create(&v)
	c.JSON(http.StatusOK, v)
}

func UpdateVote(c *gin.Context) {
	var v Vote
	uuidParam := c.Query("uuid")
	config.DB.Where("uuid = ?", uuidParam).Find(&v)
	var updatedVote Vote
	err := c.BindJSON(&updatedVote)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}


	v.StartDate = updatedVote.StartDate
	v.EndDate = updatedVote.EndDate

	config.DB.Save(&v)
	c.JSON(http.StatusOK, v)
}

func DeleteVote(c *gin.Context) {
	var v Vote
	uuidParam := c.Query("uuid")
	config.DB.Where("uuid = ?", uuidParam).Find(&v)
	config.DB.Delete(&v)
	c.JSON(http.StatusOK, v)
}
