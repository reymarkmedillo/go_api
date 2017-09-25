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

// "Before Middleware"
func BeforeMiddleware(c *gin.Context) {
	db := Database()
	var client Client
	secret := c.Param("client_secret")
	key := c.Param("client_key")

	db.Where("client_secret = ? and client_key = ?", secret, key).First(&client)
	fmt.Println(client)
}

// "Token Middleware"
func TokenMiddleware(c *gin.Context) {
	jwtString := c.GetHeader("Authorization")
	token, err := jwt.Parse(jwtString, keyFn)
	claims, ok := token.Claims.(jwt.MapClaims)

	nowDate := time.Now().Format(dateFormat)
	var tokenDate = claims["expires_at"].(string)
	cnvTokenDate, _ := time.Parse(dateFormat, tokenDate)
	cnvNowDate, _ := time.Parse(dateFormat, nowDate)
	isTokenExpired := cnvNowDate.After(cnvTokenDate)

	fmt.Printf("date now %s \n", nowDate)
	fmt.Printf("date token %s \n", tokenDate)

	if err == nil && ok && token.Valid {
		if isTokenExpired {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "token already expired"})
			return
		}
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "invalid token"})
		return
	}
}

func keyFn(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return []byte("Base"), nil
}

func userProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}

func loginHandler(c *gin.Context) {
	db := Database()
	var user User
	_pass := []byte(c.PostForm("password"))
	fmt.Println(c.PostForm("email"))
	db.Where("email = ?", c.PostForm("email")).First(&user)

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
	expiryDate := time.Now().Add(time.Minute * 5).Format(dateFormat) // add 5 minutes
	claims["expires_at"] = expiryDate

	tokenString, _ := token.SignedString(mySigningKey)
	c.JSON(http.StatusCreated, gin.H{
		"status":       http.StatusCreated,
		"access_token": tokenString,
		"expires_at":   expiryDate,
	})

	accessToken := AccessToken{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: expiryDate,
	}

	db.Save(&accessToken)
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
	Type     string `json:"auth_type"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Role     string `json:"role"`
	Valid    int    `json:"valid"`
	Remember string `json:"remember_token"`
}

// "AccessToken Model"
type AccessToken struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// "Client Model"
type Client struct {
	gorm.Model
	ClientKey    string `json:"client_key"`
	ClientSecret string `json:"client_secret"`
}
