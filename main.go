package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
	"time"
)

func task() {
	log.Println("Task is being performed.")
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

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Minutes().Do(task)
	s.StartAsync()

	router.Run(":" + port)
}
