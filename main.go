package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
	"os"
)

func task() {
	fmt.Println("Task is being performed.")
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	s := gocron.NewScheduler()
	s.Every(2).Hours().Do(task)
	<-s.Start()

	router.Run(":" + port)
}
