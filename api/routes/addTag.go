package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Abhinav-987/url-shortner/api/database"
	"github.com/Abhinav-987/url-shortner/api/models"
	"github.com/gin-gonic/gin"
)

func AddTag(c *gin.Context) {
	var tagRequest models.TagRequest
	if err := c.ShouldBind(&tagRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	shortID := tagRequest.ShortID
	tag := tagRequest.Tag

	val, err := database.Client.Get(database.Ctx, shortID).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found for the given shortID",
		})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		//if the data is not a JSON object, assume it as a plain string
		data = make(map[string]interface{})
		data["data"] = val
	}

	//check if tag fields already exist aand it's a slice of strings
	var tags []string
	if existingTags, ok := data["tags"].([]interface{}); ok {
		for _, t := range existingTags {
			if strTag, ok := t.(string); ok {
				tags = append(tags, strTag)
			}
		}
	}

	//check for duplicate tags
	for _, existingTag := range tags {
		if existingTag == tag {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Tag already exists",
			})
			return
		}
	}

	//add new tag to the tags slice
	tags = append(tags, tag)
	data["tags"] = tags

	//Marshall the updated data back to json
	updatedData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to marshall the updated data",
		})
		return
	}

	err = database.Client.Set(database.Ctx, shortID, updatedData, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update database",
		})
		return
	}

	c.JSON(http.StatusOK, data)

}
