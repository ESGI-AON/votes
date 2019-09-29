package helper

import (
	"fmt"
	"github.com/google/uuid"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)


func GetAccessLevel(claims jwt.MapClaims) int {
	return int(claims["accessLevel"].(float64))
}

func GetUUID(claims jwt.MapClaims) uuid.UUID {
	UUID, _ := uuid.Parse(fmt.Sprintf("%v", claims["uuid"]))
	return UUID
}

func GetUUIDStr(claims jwt.MapClaims) string {
	return fmt.Sprintf("%v", claims["uuid"])
}

