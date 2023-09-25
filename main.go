package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
	"github.com/joho/godotenv"

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

	router.POST("api/image", uploadImage)

	fmt.Printf("Server is running on  http://localhost:%s\n", PORT)
	if err := router.Run(":" + PORT); err != nil {
		log.Fatal(err)
	}
}

func uploadImage(c *gin.Context) {
	qualityStr := c.DefaultQuery("quality", "70")
	extension := c.DefaultQuery("extension", "webp")
	quality, err := strconv.Atoi(qualityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid 'quality' parameter. It must be a valid integer.",
		})
		return
	}

	// Check if the image format is allowed
	allowedFormat, ok := AllowedExtensions[extension]

	var supportedExtension string

	// Iterate through the map and append values to the string
	for _, value := range AllowedExtensions {
		supportedExtension += value + ","
	}

	// Remove the trailing comma, if any
	supportedExtension = strings.TrimSuffix(supportedExtension, ",")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Unsupported image extension: %s. Supported extensions are %s", allowedFormat, supportedExtension),
		})
		return
	}

	fileheader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	file, err := fileheader.Open()
	if err != nil {
		panic(err)
	}

	buffer, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	errDir := createFolder("static")
	if errDir != nil {
		panic(errDir)
	}

	filename, err := imageProcessing(buffer, quality, "static", extension)
	if err != nil {
		panic(err)
	}
	relativePath := fmt.Sprintf("http://localhost:%s/uploads/"+filename, PORT)
	c.JSON(http.StatusOK, gin.H{
		"compressed": relativePath,
	})
}

var AllowedExtensions = map[string]string{
	"jpeg": "jpeg",
	"svg":  "svg",
	"png":  "png",
	"webp": "webp",
	"heif": "heif",
	"avif": "avif",
	"gif":  "gif",
}

var extensionToImageType = map[string]bimg.ImageType{
	"jpeg": bimg.JPEG,
	"png":  bimg.PNG,
	"svg":  bimg.SVG,
	"webp": bimg.WEBP,
	"heif": bimg.HEIF,
	"avif": bimg.AVIF,
	"gif":  bimg.GIF,
}

func imageProcessing(buffer []byte, quality int, dirname string, extension string) (string, error) {

	filename := strings.Replace(uuid.New().String(), "-", "", -1) + "." + extension

	converted, err := bimg.NewImage(buffer).Convert(extensionToImageType[extension])
	if err != nil {
		return filename, err
	}

	compressed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
	if err != nil {
		return filename, err
	}

	writeError := bimg.Write(fmt.Sprintf("./"+dirname+"/%s", filename), compressed)
	if writeError != nil {
		return filename, writeError
	}

	return filename, nil
}

func createFolder(dirname string) error {
	_, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirname, 0755)
		if errDir != nil {
			return errDir
		}
	}
	return nil
}
