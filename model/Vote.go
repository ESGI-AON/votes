package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"strings"
	"time"
)


type Vote struct {
	ID        uint `gorm:"primary_key" json:"id"`
	UUID uuid.UUID `gorm:"not null" json:"uuid"`
	Title     string `json:"title"`
	Description string `json:"desc"`
	UUIDVote    pq.StringArray  `gorm:"type:varchar(100)[]" json:"uuid_vote"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `json:"up	dated_at"`
	DeletedAt *time.Time `json:"deleted_at" pg:",soft_delete"`
}

type VoteResponse struct {
	UUID uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"desc"`
}

// set title
func (vo *Vote) SetTitle(title string) {
	if title != "" {
		vo.Title = title
	}
}

// set description
func (vo *Vote) SetDescription(desc string) {
	if desc != "" {
		vo.Description = desc
	}
}

// set uuid and created_at before creating resource
func (vo *Vote) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UUID", uuid.New())
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

func (vo *Vote) SetStartDate(date time.Time) {
	defaultTime := time.Time{}
	fmt.Println(defaultTime == date)
	if date == defaultTime {
		return
	}
	vo.StartDate = date
}

func (vo *Vote) SetEndDate(date time.Time) {
	defaultTime := time.Time{}
	if date == defaultTime {
		return
	}
	vo.EndDate = date
}

// marshal struct to json
func (vo Vote) MarshalJSON() ([]byte, error) {
	var vr VoteResponse
	vr.UUID = vo.UUID
	vr.Title = vo.Title
	vr.Description = vo.Description
	//vr.StartDate = vo.StartDate.Format("02-01-2006")
	//vr.EndDate = vo.EndDate.Format("02-01-2006")

	return json.Marshal(vr)
}

// unmarshal json to struct
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
		if strings.ToLower(k) == "desc" {
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
