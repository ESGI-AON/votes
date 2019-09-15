package model

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"regexp"
	"strings"
	"time"
)


type User struct {
	ID        uint `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `gorm:"not null" json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt *time.Time `json:"-"`
	UUID uuid.UUID `gorm:"not null" json:"uuid"`
	AccessLevel int `gorm:"not null" json:"-"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName string `gorm:"not null" json:"last_name"`
	Email string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	DateOfBirth time.Time `json:"-"`
}

type UserResponse struct {
	UUID uuid.UUID `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email  string `json:"email"`
}

func (u User) IsValid() []error{
	var errs []error
	firstname := strings.Trim(u.FirstName, " ")
	lastname := strings.Trim(u.LastName, " ")
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(u.Email) {
		errs = append(errs, errors.New("Email address is not valid"))
	}
	if strings.Contains(firstname, " "){
		errs = append(errs, errors.New("FirstName can't have spaces"))
	}
	if strings.Contains(lastname, " "){
		errs = append(errs, errors.New("LastName can't have spaces"))
	}
	if len(firstname) < 2 {
		errs = append(errs, errors.New("FirstName must be at least 2 characters"))
	}
	if len(lastname) < 2 {
		errs = append(errs, errors.New("LastName must be at least 2 characters"))
	}
	// TODO compare dates to check if > 18
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (u *User) SetPassword(pwd string) {
	h := sha256.New()
	h.Write([]byte(pwd))
	u.Password = base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UUID", uuid.New())
	return nil
}

func (u User) MarshalJSON() ([]byte, error) {
	var ur UserResponse
	ur.UUID = u.UUID
	ur.FirstName = u.FirstName
	ur.LastName = u.LastName
	ur.Email = u.Email
	return json.Marshal(ur)
}
