package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nirajpdn/image-compress/src/util"
)

func UploadImage(c *gin.Context) {
	PORT := os.Getenv("PORT")
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
	allowedFormat, ok := util.AllowedExtensions[extension]

	var supportedExtension string

	// Iterate through the map and append values to the string
	for _, value := range util.AllowedExtensions {
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

	errDir := util.CreateFolder("static")
	if errDir != nil {
		panic(errDir)
	}

	filename, err := util.ImageProcessing(buffer, quality, "static", extension)
	if err != nil {
		panic(err)
	}
	relativePath := fmt.Sprintf("http://localhost:%s/uploads/"+filename, PORT)
	c.JSON(http.StatusOK, gin.H{
		"compressed": relativePath,
	})
}
