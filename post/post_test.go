package post_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"example.com/social-gin/post"
	"example.com/social-gin/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var postHandler *post.Handler
var postId string

func TestMain(m *testing.M) {
	dsn := "sqlserver://sert11:1234567890@localhost:1433?database=social1"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&post.Post{})
	// prepare handler
	postHandler = &post.Handler{
		DB: db,
	}

	os.Exit(m.Run())
}

func TestAddPostCaseStatusOk(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()

	r.POST("/users/:uid/posts", postHandler.AddPost)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Content": "test post of userid 1",
		"Likes":   1,
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPost, "/users/1/posts", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnPost := post.Post{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnPost)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test post of userid 1"
	get := returnPost.Content
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "1"
	get = strconv.Itoa(returnPost.Likes)
	status = assert.Equal(t, want, get)
	if !status {
		return
	}
	postId = strconv.Itoa(int(returnPost.ID))

}

func TestUpdatePostCaseStatusOk(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()

	r.PUT("/users/:uid/posts/:pid", postHandler.UpdatePost)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Content": "test Update post of userid 1",
		"Likes":   2,
	})
	given := string(givenBytes)
	updatePostUrl := "/users/1/posts/" + postId
	req := httptest.NewRequest(http.MethodPut, updatePostUrl, strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnPost := post.Post{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnPost)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test Update post of userid 1"
	get := returnPost.Content
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "2"
	get = strconv.Itoa(returnPost.Likes)
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

}

func TestGetPostCaseStatusOk(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()

	r.GET("/users/:uid/posts/:pid", postHandler.GetPost)

	updatePostUrl := "/users/1/posts/" + postId
	req := httptest.NewRequest(http.MethodGet, updatePostUrl, nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnPost := post.Post{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnPost)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test Update post of userid 1"
	get := returnPost.Content
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "2"
	get = strconv.Itoa(returnPost.Likes)
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

}

func TestDeletePostCaseStatusOk(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()

	r.DELETE("/users/:uid/posts/:pid", postHandler.GetPost)

	updatePostUrl := "/users/1/posts/" + postId
	req := httptest.NewRequest(http.MethodDelete, updatePostUrl, nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnPost := post.Post{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnPost)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test Update post of userid 1"
	get := returnPost.Content
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "2"
	get = strconv.Itoa(returnPost.Likes)
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

}
