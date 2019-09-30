package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/votes/config"
	"github.com/votes/controller"
	"github.com/votes/model"
	"log"
	"net/http"
	"os"
	"time"
)


type login struct {
	Email string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type User = model.User
type Vote = model.Vote

var err error

func main(){

	config.DB, err = gorm.Open("postgres", "host=localhost port=5432 user=root password=root dbname=govotes sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}

	config.DB.AutoMigrate(&User{} , &Vote{})

	port := os.Getenv("PORT")
	r := gin.Default()

	if port == "" {
		port = "8000"
	}

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
			return &user,nil

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			claims := jwt.ExtractClaims(c)
			if len(claims) > 0 {
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
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc: time.Now,
	}


	// AUTH
	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/")
	auth.POST("users", controller.CreateUser)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		// USER
		auth.GET("users/:uuid", controller.GetUser)
		auth.PUT("users/:uuid", controller.UpdateUser)
		auth.DELETE("users/:uuid", controller.DeleteUser)
		// VOTES
		auth.GET("votes/:uuid", controller.GetVote)
		auth.POST("votes", controller.CreateVote)
		auth.PUT("votes/:uuid", controller.UpdateVote)
	}


	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})


	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}