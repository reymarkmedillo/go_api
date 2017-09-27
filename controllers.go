package main

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func userProfile(c *gin.Context) {
	jwtString := c.GetHeader("Authorization")
	token, _ := jwt.Parse(jwtString, keyFn)
	claims, _ := token.Claims.(jwt.MapClaims)

	db := Database()
	type Result struct {
		Email         string
		FirstName     string
		LastName      string
		Address       string
		Premium       int
		PaymentMethod string
	}
	var result Result
	db.Table("users").Select("users.email, user_profiles.first_name, user_profiles.last_name,user_profiles.address,user_profiles.premium,user_profiles.payment_method").Joins("LEFT JOIN user_profiles on user_profiles.user_id = users.id").Where("email = ?", claims["email"]).First(&result)

	c.JSON(http.StatusOK, gin.H{"user": result})

}

func loginHandler(c *gin.Context) {
	db := Database()
	var user User
	_pass := []byte(c.PostForm("password"))
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

	claims["role"] = user.Role
	claims["id"] = user.ID
	claims["email"] = c.PostForm("email")
	expiryDate := time.Now().Add(time.Hour * 5).Format(dateFormat) // add 5 hours
	claims["expires_at"] = expiryDate

	tokenString, _ := token.SignedString(mySigningKey)
	c.JSON(http.StatusCreated, gin.H{
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
