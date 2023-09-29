package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nirajpdn/image-compress/src/controller"

	"github.com/gin-gonic/gin"
)

var PORT string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "4000"
	}
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hi there. I'm up.",
		})
	})

	router.StaticFS("/uploads", http.Dir("static"))

	router.POST("api/image", controller.UploadImage)

	fmt.Printf("Server is running on  http://localhost:%s\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		log.Fatal(err)
	}
}
