package main

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func searchCases(c *gin.Context) {
	db := Database()
	grs := []Case{}
	title := []Case{}
	syllabus := []Case{}
	topic := []Case{}

	var buffer bytes.Buffer
	search := c.PostForm("search")

	buffer.WriteString("%")
	buffer.WriteString(search)
	buffer.WriteString("%")

	// SEARCH IN GR NUMBER
	db.Table("cases").Where("grno LIKE ?", buffer.String()).Scan(&grs)
	// SEARCH IN TITLE
	db.Table("cases").Where("title LIKE ?", buffer.String()).Scan(&title)
	// SEARCH IN SYLLABUS
	db.Table("cases").Where("syllabus LIKE ?", buffer.String()).Scan(&syllabus)
	// SEARCH IN TOPIC
	db.Table("cases").Where("topic LIKE ?", buffer.String()).Scan(&topic)

	result := append(topic, syllabus...)

	c.JSON(http.StatusOK, gin.H{"result": result})
}
