package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func main() {

	// setup configuration
	viper.SetDefault("port", ":1323")
	viper.SetDefault("dsn", "sqlserver://sert11:1234567890@localhost:1433?database=social1")
	viper.SetDefault("redis", "localhost:1433")
	viper.SetDefault("redispass", "GoLang")
	viper.SetDefault("redisdb", 0)
	viper.AutomaticEnv()

	// prepare logger
	l, _ := zap.NewProduction()
	defer l.Sync()

	// prepare database
	dsn := viper.GetString("dsn")
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// setup connection pool
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDb.SetMaxIdleConns(1)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxIdleTime(time.Minute)
	sqlDb.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&user.User{}, &post.Post{}); err != nil {
		log.Fatal(err)
	}

	// prepare redis client

	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis"),
		Password: viper.GetString("redispass"),
		DB:       viper.GetInt("redisdb"),
	})
	var ctx = context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatal(err)
	}

	// prepare handler
	userHandler := &user.Handler{
		DB: db,
	}
	postHandler := &post.Handler{
		DB: db,
	}
	authHandler := &auth.Handler{
		DB:          db,
		RedisClient: client,
	}

	// prepare router
	r := gin.New()
	r.Use(gin.Recovery())
	// r.Use(gin.Logger())

	r.Use(logger.Middleware(l))

	// Routes

	r.GET("/hello", userHandler.Hello)

	r.POST("/login", authHandler.LogIn)

	r.GET("/users", userHandler.ListUser)
	r.GET("/users/:uid", userHandler.GetUser)
	r.POST("/users", userHandler.AddUser)

	r.GET("/users/:uid/posts", postHandler.ListPost)
	r.GET("/users/:uid/posts/:pid", postHandler.GetPost)

	// authenticated group
	g := r.Group("", authHandler.Authorize)

	g.GET("/tables/:uid", userHandler.ListTable)
	g.PUT("/users/:uid", userHandler.UpdateUser)
	g.DELETE("/users/:uid", userHandler.DeleteUser)

	g.POST("/users/:uid/posts", postHandler.AddPost)
	g.PUT("/users/:uid/posts/:pid", postHandler.UpdatePost)
	g.DELETE("/users/:uid/posts/:pid", postHandler.DeletePost)
	// start server
	srv := &http.Server{
		Addr:    viper.GetString("port"),
		Handler: r,
	}

	go func() {
		log.Println("Starting server at", viper.GetString("port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
