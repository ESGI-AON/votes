package model

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"regexp"
	"strings"
	"time"
)


type User struct {
	gorm.Model
	UUID uuid.UUID `gorm:"not null"`
	AccessLevel int `gorm:"not null"`
	FirstName string `gorm:"not null"`
	LastName string `gorm:"not null"`
	Email string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	DateOfBirth time.Time
}

func (user User) IsValid() []error{
	var errs []error
	firstname := strings.Trim(user.FirstName, " ")
	lastname := strings.Trim(user.LastName, " ")
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !re.MatchString(user.Email) {
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
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UUID", uuid.New())
	return nil
}
