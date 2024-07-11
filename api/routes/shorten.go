package routes

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Abhinav-987/url-shortner/api/database"
	"github.com/Abhinav-987/url-shortner/api/models"
	"github.com/Abhinav-987/url-shortner/api/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func ShortenURL(c *gin.Context) {
	var body models.Request

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Parse JSON"})
		return
	}

	clientIP := c.ClientIP()

	// Log the initial API quota and ClientIP
	apiQuota := os.Getenv("API_QUOTA")
	log.Printf("API_QUOTA: %s, ClientIP: %s", apiQuota, clientIP)

	val, err := database.Client.Get(database.Ctx, clientIP).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Error getting rate limit for %s: %v", clientIP, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	log.Printf("Rate limit value for %s: %s", clientIP, val)

	if err == redis.Nil {
		if apiQuota == "" {
			apiQuota = "10" // Default to 10 if not set
		}
		log.Printf("Setting initial quota for %s: %s", clientIP, apiQuota)
		err = database.Client.Set(database.Ctx, clientIP, apiQuota, 30*60*time.Second).Err()
		if err != nil {
			log.Printf("Error setting initial quota for %s: %v", clientIP, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		val = apiQuota
	}

	valInt, _ := strconv.Atoi(val)
	log.Printf("Current rate limit for %s: %d", clientIP, valInt)

	if valInt <= 0 {
		limit, _ := database.Client.TTL(database.Ctx, clientIP).Result()
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":            "rate limit exceeded",
			"rate_limit_reset": limit / time.Nanosecond / time.Minute,
		})
		return
	}

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	if !utils.IsDifferentDomain(body.URL) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "You can't hack this system :)",
		})
		return
	}

	body.URL = utils.EnsureHttpPrefix(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	val, _ = database.Client.Get(database.Ctx, id).Result()
	if val != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "URL Custom Short Already Exists",
		})
		return
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = database.Client.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to connect to the redis server",
		})
		return
	}

	resp := models.Response{
		Expiry:          body.Expiry,
		XRateLimitReset: 30,
		XRateRemaining:  10,
		URL:             body.URL,
		CustomShort:     "",
	}

	database.Client.Decr(database.Ctx, c.ClientIP())

	val, _ = database.Client.Get(database.Ctx, c.ClientIP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := database.Client.TTL(database.Ctx, c.ClientIP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	c.JSON(http.StatusOK, resp)
}
