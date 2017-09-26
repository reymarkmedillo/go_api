package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	dateFormat = "2006-01-02 15:04:05"
)

var mySigningKey = []byte("Base")

// "Database connection..."
func Database() *gorm.DB {
	// open connection
	db, err := gorm.Open("mysql", "root:mysql@/gotest?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database..")
	}
	return db
}

func main() {
	db := Database()
	db.AutoMigrate(&Client{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&AccessToken{})

	router := gin.Default()

	router.POST("/api/v1/auth/login", loginHandler)
	router.POST("/api/v1/auth/register", userRegistration)

	v1 := router.Group("/api/v1")
	v1.Use(BeforeMiddleware, TokenMiddleware)
	{
		v1.GET("/user", userProfile)
	}
	router.Run()
}
