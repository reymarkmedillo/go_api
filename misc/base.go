package main

import (
    "bytes"
    "database/sql"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/gotest")
    if err != nil {
        fmt.Print(err.Error())
    }
    defer db.Close()
    // make sure connection is available
    err = db.Ping()
    if err != nil {
        fmt.Print(err.Error())
    }
    type Person struct {
        Id         int
        First_Name string `form:"first_name" json:"first_name" binding:"required"`
        Last_Name  string `form:"last_name" json:"last_name" binding:"required"`
    }
    router := gin.Default()

    // GET a person detail
    router.GET("/person/:id", func(c *gin.Context) {
        var (
            person Person
            result gin.H
        )
        id := c.Param("id")
        row := db.QueryRow("select id, first_name, last_name from persons where id = ?;", id)
        err = row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
        if err != nil {
            // If no results send null
            result = gin.H{
                "result": nil,
                "count":  0,
            }
        } else {
            result = gin.H{
                "result": person,
                "count":  1,
            }
        }
        c.JSON(http.StatusOK, result)
    })

    // GET all persons
    router.GET("/persons", func(c *gin.Context) {
        var (
            person  Person
            persons []Person
        )
        rows, err := db.Query("select id, first_name, last_name from persons;")
        if err != nil {
            fmt.Print(err.Error())
        }
        for rows.Next() {
            err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
            persons = append(persons, person)
            if err != nil {
                fmt.Print(err.Error())
            }
        }
        defer rows.Close()
        c.JSON(http.StatusOK, gin.H{
            "result": persons,
            "count":  len(persons),
        })
    })

    // POST new person details
    router.POST("/person", func(c *gin.Context) {
        var (
            buffer bytes.Buffer
            person Person
        )
        c.Bind(&person)
        if person.First_Name == "" {
            c.JSON(http.StatusOK, gin.H{
                "message": "first name is required",
            })
            return
        }
        if person.Last_Name == "" {
            c.JSON(http.StatusOK, gin.H{
                "message": "last name is required",
            })
            return
        }

        stmt, err := db.Prepare("insert into persons (first_name, last_name) values(?,?);")
        if err != nil {
            fmt.Print(err.Error())
        }
        _, err = stmt.Exec(person.First_Name, person.Last_Name)

        if err != nil {
            fmt.Print(err.Error())
        }

        // Fastest way to append strings
        buffer.WriteString(person.First_Name)
        buffer.WriteString(" ")
        buffer.WriteString(person.Last_Name)
        defer stmt.Close()
        name := buffer.String()
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf(" %s successfully created", name),
        })
    })

    // PUT - update a person details
    router.PUT("/person", func(c *gin.Context) {
        var buffer bytes.Buffer
        id := c.Query("id")
        first_name := c.PostForm("first_name")
        last_name := c.PostForm("last_name")
        stmt, err := db.Prepare("update persons set first_name= ?, last_name= ? where id= ?;")
        if err != nil {
            fmt.Print(err.Error())
        }
        _, err = stmt.Exec(first_name, last_name, id)
        if err != nil {
            fmt.Print(err.Error())
        }

        // Fastest way to append strings
        buffer.WriteString(first_name)
        buffer.WriteString(" ")
        buffer.WriteString(last_name)
        defer stmt.Close()
        name := buffer.String()
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("Successfully updated to %s", name),
        })
    })

    // Delete resources
    router.DELETE("/person", func(c *gin.Context) {
        id := c.Query("id")
        stmt, err := db.Prepare("delete from persons where id= ?;")
        if err != nil {
            fmt.Print(err.Error())
        }
        _, err = stmt.Exec(id)
        if err != nil {
            fmt.Print(err.Error())
        }
        c.JSON(http.StatusOK, gin.H{
            "message": fmt.Sprintf("Successfully deleted user: %s", id),
        })
    })
    router.Run(":4000")
}




