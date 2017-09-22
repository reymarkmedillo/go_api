package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	ValidEmail = "asdf@email.com"
	ValidPass  = "1234"
	SecretKey  = "WOW,MuchShibe,ToDogge"
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
	db.AutoMigrate(&User{})
	db.AutoMigrate(&AccessToken{})

	router := gin.Default()

	router.POST("/api/v1/auth/login", loginHandler)
	router.POST("/api/v1/auth/register", userRegistration)

	v1 := router.Group("/api/v1")
	v1.Use(TokenMiddleware)
	{
		v1.GET("/user", userProfile)
	}

	router.Run()
}

func userProfile(c *gin.Context) {

}

func TokenMiddleware(c *gin.Context) {
	fmt.Println("dummy")
	c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"msg": "error"})
	// c.Next()
}

func loginHandler(c *gin.Context) {
	db := Database()
	var user User
	_pass := []byte(c.PostForm("password"))
	fmt.Println(c.PostForm("email"))
	db.Where("email=?", c.PostForm("email")).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"msg": "Incorrect email"})
		return
	}

	checkPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), _pass)
	if checkPassword != nil {
		c.JSON(http.StatusForbidden, gin.H{"msg": "Incorrect password"})
		return
	}

	// Create the Claims
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["role"] = "admin"
	claims["id"] = "13121"
	claims["email"] = ValidEmail
	expiryDate := time.Now().Add(time.Minute * 5).Format(dateFormat)
	claims["expires_at"] = expiryDate

	tokenString, _ := token.SignedString(mySigningKey)
	c.JSON(http.StatusCreated, gin.H{
		"status":       http.StatusCreated,
		"access_token": tokenString,
		"expires_at":   expiryDate,
	})

	// accessToken := AccessToken{
	// 	Token:      tokenString,
	// 	Expires_At: expiryDate,
	// 	Created_At: time.Now(),
	// }

	// db := Database()
	// db.Save(&accessToken)
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

// "User Model"
type User struct {
	gorm.Model
	ID       int    `form:"userid" json:"userid" binding:"required"`
	Type     string `json:"auth_type"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Role     string `json:"role"`
	Valid    int    `json:"valid"`
	Remember string `json:"remember_token"`
}

type AccessToken struct {
	ID         int    `json:"id"`
	User_Id    int    `json:"user_id"`
	Token      string `json:"token"`
	Expires_At string `json:"expires_at"`
	Created_At time.Time
}
