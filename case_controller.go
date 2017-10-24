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
	tempResult := []Case{}
	casegroup := []Case{}

	result := []Case{}

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
	// SEARCH IN CASE GROUP
	db.Table("case_groups").Where("refno LIKE ?", buffer.String()).Joins("left join cases as c on c.id = case_groups.case_id").Select("case_groups.title, case_id as id, refno as grno, c.scra, c.date, c.topic, c.syllabus,c.body,c.status").Scan(&casegroup)

	tempResult = append(tempResult, syllabus...)
	tempResult = append(tempResult, title...)
	tempResult = append(tempResult, topic...)
	tempResult = append(tempResult, grs...)
	tempResult = append(tempResult, casegroup...)

	encountered := map[uint]bool{}

	for _, v := range tempResult {
		if encountered[v.ID] == true {

		} else {
			encountered[v.ID] = true
			result = append(result, v)
		}
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func viewCase(c *gin.Context) {
	db := Database()
	var caseResult Case

	id := c.Param("case_id")
	db.Table("cases").Where("id = ?", id).First(&caseResult)

	if caseResult.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"result": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": caseResult})
}

func createDraftCase(c *gin.Context) {
	db := Database()
	newcase := Case{}
}
