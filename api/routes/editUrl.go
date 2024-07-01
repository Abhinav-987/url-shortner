package routes

import (
	"net/http"
	"time"

	"github.com/Abhinav-987/url-shortner/api/database"
	"github.com/Abhinav-987/url-shortner/api/models"
	"github.com/gin-gonic/gin"
)

func EditURL(c *gin.Context) {
	shortID := c.Param("shortID")
	var body models.Request

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot Parse JSON",
		})
		return
	}

	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil || val == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ShortID doesn't exists",
		})
		return
	}

	err = database.Client.Set(database.Ctx, shortID, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update the shortened content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The content has been updated !!!!!",
	})

}
