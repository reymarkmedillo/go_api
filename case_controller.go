package main

import (
	"bytes"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func searchCases(c *gin.Context) {
	db := Database()
	grs := []CaseResult{}
	title := []CaseResult{}
	syllabus := []CaseResult{}
	topic := []CaseResult{}
	tempResult := []CaseResult{}
	casegroup := []CaseResult{}

	result := []CaseResult{}

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
	db.Table("case_groups").Where("refno LIKE ?", buffer.String()).Joins("left join cases as c on c.id = case_groups.case_id").Select("case_groups.title, c.id, refno as grno, c.scra, c.date, c.topic, c.syllabus,c.body,c.status").Scan(&casegroup)

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

	result = makeCase(result)

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
	// db := Database()
	// newcase := Case{}
}

func makeCase(caseRes []CaseResult) []CaseResult {
	db := Database()
	childs := []Children{}
	result := []CaseResult{}

	for _, v := range caseRes {
		db.Table("case_groups").Where("case_id = ?", v.ID).Scan(&childs)
		v.Child = append(v.Child, childs...)
		result = append(result, v)
	}

	return result
}
