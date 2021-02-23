package user_test

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

var userHandler *user.Handler

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
	userHandler = &user.Handler{
		DB: db,
	}

	os.Exit(m.Run())
}

func TestAddUserCase_UserName_Nill(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.POST("/users", userHandler.AddUser)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Password": "test_pass",
		"Name":     "test names",
		"Email":    "test@example",
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code == http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusBadRequest, rec.Code)
	}

}

func TestAddUserCase_Password_Nill(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.POST("/users", userHandler.AddUser)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Username": "testAddUser",
		"Name":     "test names",
		"Email":    "test@example",
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code == http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusBadRequest, rec.Code)
	}

}

func TestAddUserCaseStatusOk(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.POST("/users", userHandler.AddUser)
	r.DELETE("/users/:uid", userHandler.DeleteUser)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Username": "testAddUser",
		"Password": "test_pass",
		"Name":     "test names",
		"Email":    "test@example",
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnUser := user.User{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnUser)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "testAddUser"
	get := returnUser.Username
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "test names"
	get = returnUser.Name

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test@example"
	get = returnUser.Email

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	//Delete User this Test Add Data
	deleteUrl := "/users/" + strconv.Itoa(int(returnUser.ID))
	req = httptest.NewRequest(http.MethodDelete, deleteUrl, strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

}

func TestUpdateUserCaseStatusOk(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.PUT("/users/:uid", userHandler.UpdateUser)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Username": "testUpdateUser1",
		"Password": "testUpdatePass1",
		"Name":     "test names Update 1",
		"Email":    "test1@example",
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnUser := user.User{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnUser)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "testUpdateUser1"
	get := returnUser.Username
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "testUpdatePass1"
	get = returnUser.Password
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "test names Update 1"
	get = returnUser.Name

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test1@example"
	get = returnUser.Email

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

}

func TestUpdateUserCaseUserName(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.PUT("/users/:uid", userHandler.UpdateUser)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Username": "test1",
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnUser := user.User{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnUser)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test1"
	get := returnUser.Username
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "testUpdatePass1"
	get = returnUser.Password
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "test names Update 1"
	get = returnUser.Name

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test1@example"
	get = returnUser.Email

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

}

func TestUpdateUserCasePassword(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.PUT("/users/:uid", userHandler.UpdateUser)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Password": "1234567890",
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnUser := user.User{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnUser)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test1"
	get := returnUser.Username
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "1234567890"
	get = returnUser.Password
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "test names Update 1"
	get = returnUser.Name

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test1@example"
	get = returnUser.Email

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

}

func TestUpdateUserCaseName(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.PUT("/users/:uid", userHandler.UpdateUser)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Name": "test names",
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnUser := user.User{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnUser)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test1"
	get := returnUser.Username
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "1234567890"
	get = returnUser.Password
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "test names"
	get = returnUser.Name

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test1@example"
	get = returnUser.Email

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

}

func TestUpdateUserCaseEmail(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.PUT("/users/:uid", userHandler.UpdateUser)

	givenBytes, _ := json.Marshal(map[string]interface{}{
		"Email": "test@example",
	})
	given := string(givenBytes)
	req := httptest.NewRequest(http.MethodPut, "/users/1", strings.NewReader(given))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnUser := user.User{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnUser)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test1"
	get := returnUser.Username
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "1234567890"
	get = returnUser.Password
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "test names"
	get = returnUser.Name

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test@example"
	get = returnUser.Email

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

}

func TestGetUserCaseByID(t *testing.T) {

	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// register your routes
	r := gin.Default()
	r.GET("/users/:uid", userHandler.GetUser)
	given := "1"
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(rec, req)

	// Check to see if the response was what you expected
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, rec.Code)
	}

	returnUser := user.User{}

	err := json.Unmarshal(rec.Body.Bytes(), &returnUser)
	if err != nil {
		t.Error("can not Unmarshal respond", string(rec.Body.Bytes()))
		return
	}

	want := "test1"
	get := returnUser.Username
	status := assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "1234567890"
	get = returnUser.Password
	status = assert.Equal(t, want, get)
	if !status {
		return
	}

	want = "test names"
	get = returnUser.Name

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

	want = "test@example"
	get = returnUser.Email

	if get != want {
		t.Error("given", given, "want", want, "but get", get)
	}

}
