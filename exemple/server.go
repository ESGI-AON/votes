package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/votes/config"
	"github.com/votes/controller"
	"github.com/votes/model"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

type login struct {
	Email string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

// User demo
type User = model.User
var err error

func main() {

	config.DB, err = gorm.Open("postgres", "host=localhost port=5432 user=root password=root dbname=govotes sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}

	config.DB.AutoMigrate(&User{})


	port := os.Getenv("PORT")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if port == "" {
		port = "8000"
	}
	// the jwt middleware
	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"uuid": v.UUID,
					"accessLevel": v.AccessLevel,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			var user model.User
			if err := c.Bind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Email
			password := loginVals.Password

			h := sha256.New()
			h.Write([]byte(password))
			password = base64.URLEncoding.EncodeToString(h.Sum(nil))

			if err := config.DB.Where("email = ? AND password = ?", userID, password).Find(&user).Error; err != nil {
				return "", jwt.ErrFailedAuthentication
			}
			return &User{Email:user.Email, AccessLevel:user.AccessLevel},nil

		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			claims := jwt.ExtractClaims(c)
			fmt.Println(claims["accessLevel"], reflect.TypeOf(claims["accessLevel"]))
			if claims["accessLevel"] == 1.00 {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value.
		// This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}


	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/user", controller.CreateUser)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
