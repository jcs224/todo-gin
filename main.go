package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Todo struct {
	gorm.Model
	Text      string `json: "text"`
	Completed bool   `json: "completed"`
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

	api := r.Group("/api")
	{
		api.POST("/todos", func(c *gin.Context) {
			todoText := c.PostForm("text")

			newTodo := Todo{Text: todoText, Completed: false}
			db.Create(&newTodo)

			fmt.Println(todoText)
		})

		api.GET("/todos", func(c *gin.Context) {
			var todos []Todo
			db.Find(&todos)

			todosJSONBytes, err := json.Marshal(todos)
			if err != nil {
				log.Fatal("Cannot encode to json", err)
			}

			todosJSON := string(todosJSONBytes)

			c.JSON(http.StatusOK, todosJSON)
		})

		api.DELETE("/todos/:id", func(c *gin.Context) {
			var todo Todo
			db.Where("id = ?", c.Param("id")).First(&todo)
			db.Unscoped().Delete(&todo)
			c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
		})
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Run()
}
