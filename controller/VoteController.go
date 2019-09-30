package controller

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/votes/config"
	"github.com/votes/helper"
	"github.com/votes/model"
	"log"
	"net/http"
	"time"
)

type GetResponse struct {
	UUID        uuid.UUID      `json:"uuid,string"`
	Title       string         `json:"title"`
	Description string         `json:"desc"`
	UUIDVote    pq.StringArray `json:"uuid_votes"`
}

type PutResponse struct {
	UUID        uuid.UUID      `json:"uuid,string"`
	Title       string         `json:"title"`
	Description string         `json:"desc"`
	StartDate   time.Time      `json:"start_date,string"`
	EndDate     time.Time      `json:"end_date,string"`
	UUIDVote    pq.StringArray `json:"uuid_votes"`
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

type Vote = model.Vote

func GetVote(c *gin.Context) {
	var vote Vote
	uuidParam := c.Param("uuid")
	config.DB.Where("uuid = ?", uuidParam).Find(&vote)
	c.JSON(http.StatusOK, GetResponse{
		UUID:        vote.UUID,
		Title:       vote.Title,
		Description: vote.Description,
		UUIDVote:    vote.UUIDVote,
	})
}

func CreateVote(c *gin.Context) {
	var v Vote
	err := c.BindJSON(&v)
	claims := jwt.ExtractClaims(c)
	accessLevel := helper.GetAccessLevel(claims)
	if accessLevel == 0{
		c.JSON(http.StatusUnauthorized, "You need to be an admin to create a vote")
		return
	}
	if v.Description == "" || v.Title == ""{
		c.JSON(http.StatusBadRequest, "You need to fill Title and description")
		return
	}
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	config.DB.NewRecord(v)
	config.DB.Create(&v)
	c.JSON(http.StatusOK, v)
}

func UpdateVote(c *gin.Context, ) {
	var v Vote
	uuidParam := c.Param("uuid")
	claims := jwt.ExtractClaims(c)
	accessLevel := helper.GetAccessLevel(claims)
	voterUUID := helper.GetUUIDStr(claims)
	fmt.Println(claims, voterUUID)

	config.DB.Where("uuid = ?", uuidParam).Find(&v)
	var updatedVote Vote
	err := c.BindJSON(&updatedVote)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if !Contains(v.UUIDVote, voterUUID) {
		v.UUIDVote = append(v.UUIDVote, voterUUID)
	}
	if accessLevel == 1 {
		v.SetTitle(updatedVote.Title)
		v.SetDescription(updatedVote.Description)
	}
	v.StartDate = updatedVote.StartDate
	v.EndDate = updatedVote.EndDate

	config.DB.Save(&v)
	c.JSON(http.StatusOK, PutResponse{
		UUID:        v.UUID,
		Title:       v.Title,
		Description: v.Description,
		StartDate:   v.StartDate,
		EndDate:     v.EndDate,
		UUIDVote:    v.UUIDVote,
	})
}
