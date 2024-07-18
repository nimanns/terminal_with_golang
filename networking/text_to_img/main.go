package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/generate_image", generate_image)

	r.Run(":8080")
}

func generate_image(c *gin.Context) {
	var request struct {
		Text string `json:"text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dc := gg.NewContext(400, 200)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFace("./Atop-R99O3.ttf", 32); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load font"})
		return
	}

	dc.DrawStringAnchored(request.Text, 200, 100, 0.5, 0.5)

	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode image"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=text_image.png"))
	c.Data(http.StatusOK, "image/png", buf.Bytes())
}
