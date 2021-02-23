package user

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"example.com/social-gin/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// User represents user data
type User struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"update_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	Username  string       `gorm:"uniquekey" json:"username"`
	Password  string       `json:"password"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
}

// Handler represents handler of user data
type Handler struct {
	DB *gorm.DB
}

// Hello handles hello request
func (h *Handler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}

// ListTable handles list table request
func (h *Handler) ListTable(c *gin.Context) {

	// custom logger
	logger := logger.Extract(c)
	uid := c.GetString("uid")
	logger.Info("listing table", zap.String("uid", uid))

	rows, err := h.DB.Raw("sp_tables").Rows()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tables := []string{}
	var tableQualifier sql.NullString
	var tableOwner sql.NullString
	var tableName sql.NullString
	var tableType sql.NullString
	var remarks sql.NullString
	for rows.Next() {

		if err := rows.Scan(&tableQualifier, &tableOwner, &tableName, &tableType, &remarks); err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		tables = append(tables, tableName.String)
	}
	c.JSON(http.StatusOK, tables)
}

// AddUser handle add user request
func (h *Handler) AddUser(c *gin.Context) {
	user := User{}
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "username is empty",
		})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "password is empty",
		})
		return
	}

	if result := h.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ListUser handle list user request
func (h *Handler) ListUser(c *gin.Context) {
	users := []User{}
	if result := h.DB.Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser handle list user request
func (h *Handler) GetUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	user := User{}
	if result := h.DB.First(&user, uid); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser handle update user request
func (h *Handler) UpdateUser(c *gin.Context) {

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	user := User{}
	if result := h.DB.First(&user, uid); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}

	updateUser := User{}
	if err := c.Bind(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// update fields
	if updateUser.Username != "" {
		user.Username = updateUser.Username
	}
	if updateUser.Password != "" {
		user.Password = updateUser.Password
	}
	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	if result := h.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser handle delete user request
func (h *Handler) DeleteUser(c *gin.Context) {

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	user := User{}
	if result := h.DB.First(&user, uid); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}

	if result := h.DB.Delete(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}
