package main

import "github.com/jinzhu/gorm"

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
