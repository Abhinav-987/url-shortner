package routes

import (
	"net/http"

	"github.com/Abhinav-987/url-shortner/api/database"
	"github.com/gin-gonic/gin"
)

func GetByShortID(c *gin.Context) {
	shortID := c.Param("shortID")

	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found for given short ID",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": val,
	})
}
