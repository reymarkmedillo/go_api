package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// "Before Middleware"
func BeforeMiddleware(c *gin.Context) {
	db := Database()
	var client Client
	secret := c.Query("client_secret")
	key := c.Query("client_key")

	db.Where("client_secret = ? and client_key = ?", secret, key).First(&client)
	if client.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Unauthorized"})
		return
	}
	c.Next()
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
