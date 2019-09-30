package helper

import (
	"fmt"
	"github.com/google/uuid"
	jwt "gopkg.in/dgrijalva/jwt-go.v3"
)

// get access level from jwt
func GetAccessLevel(claims jwt.MapClaims) int {
	return int(claims["accessLevel"].(float64))
}

// get uuid from jwt
func GetUUID(claims jwt.MapClaims) uuid.UUID {
	UUID, _ := uuid.Parse(fmt.Sprintf("%v", claims["uuid"]))
	return UUID
}

// get uuid as string from jwt
func GetUUIDStr(claims jwt.MapClaims) string {
	return fmt.Sprintf("%v", claims["uuid"])
}

