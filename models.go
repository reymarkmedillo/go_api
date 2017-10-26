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

type UserProfile struct {
	gorm.Model
	UserID        uint   `json:"user_id" gorm:"joins:users;on:id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Address       string `json:"address"`
	Premium       int    `json:"premium"`
	PaymentMethod string `json:"payment_method"`
}

type Case struct {
	gorm.Model
	Title    string `gorm:"type:mediumtext"`
	Scra     string `gorm:"type:text"`
	Grno     string
	Date     string `gorm:"type:date"`
	Topic    string `gorm:"type:mediumtext"`
	Syllabus string `gorm:"type:longtext"`
	Body     string `gorm:"type:longtext"`
	Status   string `gorm:"type:text"`
}

type CaseResult struct {
	gorm.Model
	Title    string
	Scra     string
	Grno     string
	Date     string
	Topic    string
	Syllabus string
	Body     string
	Status   string
	Child    []Children
}

type Children struct {
	Refno string
	Title string
}

type CaseGroup struct {
	gorm.Model
	CaseID uint   `gorm:"type:int;size:11"`
	Refno  string `gorm:"type:text"`
	Title  string `gorm:"type:mediumtext"`
}

// "MODEL FOR RESULTS ..."
type UserProfileResult struct {
	Email         string
	FirstName     string
	LastName      string
	Address       string
	Premium       int
	PaymentMethod string
}
