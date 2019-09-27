package model

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"regexp"
	"strconv"
	"strings"
	"time"
)


type User struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" pg:",soft_delete"`
	UUID uuid.UUID `gorm:"not null" json:"uuid"`
	AccessLevel int64 `gorm:"not null" json:"access_level,string"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName string `gorm:"not null" json:"last_name"`
	Email string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	DateOfBirth time.Time `json:"birth_date"`
}

type UserResponse struct {
	UUID uuid.UUID `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email  string `json:"email"`
	DateOfBirth string `json:"birth_date"`
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
	if pwd == "" {
		return
	}
	h := sha256.New()
	h.Write([]byte(pwd))
	u.Password = base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (u *User) SetFirstname(name string) {
	if name != "" {
		u.FirstName = name
	}
}

func (u *User) SetLastname(name string) {
	if name != "" {
		u.LastName = name
	}
}

func (u *User) SetEmail(email string) {
	if email != "" {
		u.Email = email
	}
}

func (u *User) SetBirthDate(date time.Time) {
	// TODO check this shit
	defaultTime := time.Time{}
	if date == defaultTime {
		return
	}
	u.DateOfBirth = date
}



func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UUID", uuid.New())
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

func (u User) MarshalJSON() ([]byte, error) {
	var ur UserResponse
	ur.UUID = u.UUID
	ur.FirstName = u.FirstName
	ur.LastName = u.LastName
	ur.Email = u.Email
	ur.DateOfBirth = u.DateOfBirth.Format("02-01-2006")
	return json.Marshal(ur)
}


func (u *User) UnmarshalJSON(data []byte) error {
	var rawStrings map[string]string
	err := json.Unmarshal(data, &rawStrings)
	if err != nil {
		return err
	}
	// TODO switch case instead of if
	for k,v := range rawStrings {
		if strings.ToLower(k) == "first_name" {
			u.FirstName = v
		}
		if strings.ToLower(k) == "last_name" {
			u.LastName = v
		}
		if strings.ToLower(k) == "email" {
			u.Email = v
		}
		if strings.ToLower(k) == "password" {
			u.SetPassword(v)
		}
		if strings.ToLower(k) == "birth_date" {
			t, err := time.Parse("02-01-2006", v)
			if err != nil {
				return err
			}
			u.DateOfBirth = t
		}
		if strings.ToLower(k) == "access_level" {
			i, _ := strconv.ParseInt(v, 10, 64)
			u.AccessLevel = i
		}
	}
	return nil
}
