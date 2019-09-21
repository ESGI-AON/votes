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
	"github.com/votes/helpers"
)


type User struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UUID uuid.UUID `gorm:"not null" json:"uuid"`
	AccessLevel int `gorm:"not null" json:"access_level"`
	FirstName string `gorm:"not null" json:"first_name"`
	LastName string `gorm:"not null" json:"last_name"`
	Email string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	DateOfBirth time.Time `json:"date_of_birth,string"`
}

type UserResponse struct {
	UUID uuid.UUID `json:"uuid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email  string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
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

func (u *User) SetBirthDate(date string) {
	u.DateOfBirth = helpers.StrToTime(date, "02-01-2006")
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
	ur.DateOfBirth = u.DateOfBirth.Format("02-01-2006")
	return json.Marshal(ur)
}

//func (u *User) UnmarshalJSON(data []byte) error {
//	type Alias User
//	aux := &struct {
//		DateOfBirth string  `json:"date_of_birth"`
//		*Alias
//	}{
//		Alias: (*Alias)(u),
//	}
//	fmt.Println(aux)
//	return errors.New("test")
//	if err := json.Unmarshal(data, &aux); err != nil {
//		return err
//	}
//	//u.LastSeen = time.Unix(aux.LastSeen, 0)
//	u.DateOfBirth, _ = time.Parse("02-01-2006", aux.DateOfBirth)
//	return nil
//}

//func (u *User) UnmarshalJSON(data []byte) error {
//	var rawStrings map[string]string
//	err := json.Unmarshal(data, &rawStrings)
//	if err != nil {
//		return err
//	}
//	for k,v := range rawStrings {
//		if strings.ToLower(k) == "dateofbirth" {
//			t, err := time.Parse(time.RFC3339, v)
//			if err != nil {
//				return err
//			}
//			u.DateOfBirth = t
//		}
//	}
//	return nil
//}
