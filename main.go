package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vidman22/craigsbaybook-api/routes"
	"log"
	"net/http"
)

func main() {
	loaded := true
	// Find and read the config file
	if err := godotenv.Load(".env"); err != nil {
		loaded = false
	}

	if loaded == true {
		log.Println("Loaded the db file addr, user, db, password")
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	search := routes.NewRouter()
	r.GET("/search", search.SearchAll)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://www.craigsbaybook.com", "https://craigsbaybook.com"},
		AllowCredentials: true,
		AllowHeaders:     []string{"authorization", "Content-Type", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowMethods:     []string{"PUT", "PATH"},
	}))
	r.SetTrustedProxies([]string{"192.168.1.2"})
	r.Run("localhost:9090")
}
