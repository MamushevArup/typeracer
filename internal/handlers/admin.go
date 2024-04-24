package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const maxMemory = 2 << 20 // 2 MB

func (h *handler) addCars(c *gin.Context) {
	err := c.Request.ParseMultipartForm(maxMemory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse multipart form", "error": err.Error()})
		return
	}

	// Retrieve the file from the form data
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to get file from form data", "error": err.Error()})
		return
	}

	// Check if the file is a PNG image
	if fileHeader.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File is not a valid PNG image"})
		return
	}

	// Pass the byte slice to the AWS S3 upload function
	err = h.service.Admin.AddAvatar(c, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
