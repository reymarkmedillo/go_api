package main

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	jwt "gopkg.in/appleboy/gin-jwt.v2"
)

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
	db.AutoMigrate(&User{})

	router := gin.Default()

	jwtMiddleware := &jwt.GinJWTMiddleware{
		Realm:         "api.io",
		Key:           []byte("1234"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour * 24,
		Authenticator: authenticate,
		PayloadFunc:   payload,
	}

	router.POST("/api/v1/auth/login", jwtMiddleware.LoginHandler)
	router.POST("/api/v1/auth/register", userRegistration)

	v1 := router.Group("/api/v1")
	v1.Use(jwtMiddleware.MiddlewareFunc())
	{
		v1.GET("/user", userProfile)
		v1.GET("/refreshToken", jwtMiddleware.RefreshHandler)
	}
	router.Run()
}

func testHandler(c *gin.Context) {
}

func userRegistration(c *gin.Context) {
	password := []byte(c.PostForm("password"))
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user := User{
		Email:    c.PostForm("email"),
		Type:     c.PostForm("auth_type"),
		Password: string(hashedPassword),
		Role:     c.PostForm("role"),
		Valid:    1,
	}

	db := Database()
	db.Save(&user)
	c.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"msg":        "successfully registered",
		"resourceId": user.ID,
	})
}

func userProfile(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.String(http.StatusOK, "id: %s\nrole: %s", claims["id"], claims["role"])
}

func authenticate(email string, password string, c *gin.Context) (string, bool) {
	// it goes without saying that you'd be going to some form
	// of persisted storage, rather than doing this
	var user User
	var _pass = []byte(password)

	db := Database()
	db.Where("email = ?", email).First(&user)

	checkPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), _pass)
	if checkPassword != nil {
		return "", false
	}

	if user.ID == 0 {
		return "", false
	}

	return user.Email, true

}

func payload(email string) map[string]interface{} {
	// in this method, you'd want to fetch some user info
	// based on their email address (which is provided once
	// they've successfully logged in).  the information
	// you set here will be available the lifetime of the
	// user's sesion
	return map[string]interface{}{
		"id":   "1231",
		"role": "user",
	}
}

// "User ..."
type User struct {
	gorm.Model
	Type     string `json:"auth_type"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Valid    int    `json:"valid"`
	Remember string `json:"remember_token"`
}
