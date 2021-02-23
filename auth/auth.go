package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"example.com/social-gin/logger"
	"example.com/social-gin/user"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Handler represents handler of authentication
type Handler struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

// UserAuthResponse To handler Login Response Result
type UserAuthResponse struct {
	Token string
}

var ctx = context.Background()

// LogIn handle login request
func (h *Handler) LogIn(c *gin.Context) {

	username := c.Request.FormValue("u")
	password := c.Request.FormValue("p")

	l := logger.Extract(c)

	l.Info(fmt.Sprint("login username=", username, " password=", password))

	user := user.User{}

	if result := h.DB.Where("Username = ?", username).Limit(1).Find(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid username or password",
		})
		return
	}

	if user.Password != password {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid username or password",
		})
		return
	}

	token := uuid.New().String()

	l.Info(fmt.Sprint("new token=", token))

	if err := h.RedisClient.Set(ctx, token, user.ID, time.Hour*3).Err(); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	userAuthResponse := UserAuthResponse{}
	userAuthResponse.Token = token
	c.JSON(http.StatusOK, userAuthResponse)
	// c.JSON(http.StatusOK, map[string]interface{}{
	// 	"token": token,s
	// })
}

// Authorize authenticate using authorization header
func (h *Handler) Authorize(c *gin.Context) {

	auth := c.GetHeader("Authorization")
	prefix := "Bearer "

	if !strings.HasPrefix(auth, prefix) {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "no authorization token found in the header",
		})
		c.Abort()
		return
	}

	token := auth[len(prefix):]

	// validate token

	var uid string
	var err error
	if uid, err = h.RedisClient.Get(ctx, token).Result(); err != nil {
		// can't connect to redis
		if err == redis.Nil {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid token",
			})
		} else {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "can't connect to redis server",
			})
		}
		c.Abort()
		return
	}

	if uid != c.Param("uid") {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user unauthenticated",
		})
		c.Abort()
		return
	}

	// set user id of the authenticated user to context
	c.Set("uid", uid)
}
