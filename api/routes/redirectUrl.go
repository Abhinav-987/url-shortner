package routes

import (
	"net/http"

	"github.com/Abhinav-987/url-shortner/api/database"
	"github.com/gin-gonic/gin"
)

func RedirectURL(c *gin.Context) {
	shortID := c.Param("shortID")
	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}
	c.Redirect(http.StatusMovedPermanently, val)

}
