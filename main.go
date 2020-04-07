package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Todo struct {
	gorm.Model
	Text      string
	Completed bool
}

func main() {
	db, err := gorm.Open("mysql", "root:password@/todo_gin?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Todo{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.POST("/api/todos", func(c *gin.Context) {
		todoText := c.PostForm("text")
		fmt.Println(todoText)
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Run()
}
