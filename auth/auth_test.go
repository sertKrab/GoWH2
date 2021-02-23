package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"example.com/social-gin/auth"
	"example.com/social-gin/logger"
	"example.com/social-gin/post"
	"example.com/social-gin/user"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var userHandler *user.Handler
var authHandler *auth.Handler
var userToken string

func TestMain(m *testing.M) {
	viper.SetDefault("port", ":1323")
	viper.SetDefault("dsn", "sqlserver://sert11:1234567890@localhost:1433?database=social1")
	viper.SetDefault("redis", "localhost:1433")
	viper.SetDefault("redispass", "GoLang")
	viper.SetDefault("redisdb", 0)
	viper.AutomaticEnv()
	dsn := viper.GetString("dsn")
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error connect database")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis"),
		Password: viper.GetString("redispass"),
		DB:       viper.GetInt("redisdb"),
	})
	var ctx = context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic("error connect redis")
	}

	// Migrate the schema
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&post.Post{})
	// prepare handler
	authHandler = &auth.Handler{
		DB:          db,
		RedisClient: client,
	}

	os.Exit(m.Run())
}

func TestLoginStatusOK(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	l, _ := zap.NewProduction()
	defer l.Sync()
	r.Use(logger.Middleware(l))
	r.POST("/login", authHandler.LogIn)

	data := url.Values{}
	data.Set("u", "sert4")
	data.Set("p", "1234567890")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}
	returnAuth := auth.UserAuthResponse{}
	err := json.Unmarshal(rec.Body.Bytes(), &returnAuth)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}
	if returnAuth.Token == "" {
		t.Error("token is nil from respond", string(rec.Body.Bytes()))
		return
	}

	userToken = returnAuth.Token

}
