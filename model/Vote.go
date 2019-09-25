package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)


type Vote struct {
	ID        uint `gorm:"primary_key" json:"id"`
	UUID uuid.UUID `gorm:"not null" json:"uuid"`
	Title     string `json:"title"`
	Description string `json:"description"`
	UUIDVote    []User  `json:"uuid_vote"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `json:"up	dated_at"`
	DeletedAt *time.Time `json:"deleted_at" pg:",soft_delete"`
}

type VoteResponse struct {
	UUID uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`

}

func (vo *Vote) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UUID", uuid.New())
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

func (vo Vote) MarshalJSON() ([]byte, error) {
	var vr VoteResponse
	vr.UUID = vo.UUID
	vr.Title = vo.Title
	vr.Description = vo.Description
	vr.StartDate = vo.StartDate.Format("02-01-2006")
	vr.EndDate = vo.EndDate.Format("02-01-2006")

	return json.Marshal(vr)
}

func (vo *Vote) UnmarshalJSON(data []byte) error {
	var rawStrings map[string]string
	err := json.Unmarshal(data, &rawStrings)
	if err != nil {
		return err
	}
	for k,v := range rawStrings {
		if strings.ToLower(k) == "title" {
			vo.Title = v
		}
		if strings.ToLower(k) == "description" {
			vo.Description = v
		}
		if strings.ToLower(k) == "start_date" {
			t, err := time.Parse("02-01-2006", v)
			if err != nil {
				return err
			}
			vo.StartDate = t
		}
		if strings.ToLower(k) == "end_date" {
			t, err := time.Parse("02-01-2006", v)
			if err != nil {
				return err
			}
			vo.EndDate = t
		}
	}
	return nil
}
