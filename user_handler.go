package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// UserHandler mengelola semua endpoint terkait user
// Menggunakan struct untuk handler memungkinkan dependency injection yang lebih baik
type UserHandler struct {
	logger *logrus.Logger
	// Di sini nanti bisa ditambahkan dependencies lain seperti database, cache, dll
}

// NewUserHandler adalah constructor function untuk membuat UserHandler baru
// Pattern ini memudahkan testing karena dependencies bisa di-mock
func NewUserHandler(logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		logger: logger,
	}
}

// Data dummy untuk simulasi database
// Dalam production, ini akan diganti dengan actual database calls
var users = make(map[string]*User)

// GetUsers mengembalikan semua users
// Endpoint: GET /api/v1/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	h.logger.Info("Fetching all users")

	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	Success(c, http.StatusOK, "Users retrieved successfully", userList)
}

// GetUserByID mengembalikan user berdasarkan ID
// Endpoint: GET /api/v1/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	h.logger.Infof("Fetching user with ID: %s", id)

	user, exists := users[id]
	if !exists {
		Error(c, http.StatusNotFound, "User not found", nil)
		return
	}

	Success(c, http.StatusOK, "User retrieved successfully", user)
}

// CreateUser membuat user baru
// Endpoint: POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	// Binding dan validasi request body
	// Gin akan otomatis validasi berdasarkan struct tags
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Validation error: %v", err)
		Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Generate UUID untuk ID user baru
	// UUID lebih aman daripada auto-increment integer karena tidak predictable
	id := uuid.New().String()
	now := time.Now()

	user := &User{
		ID:        id,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	users[id] = user
	h.logger.Infof("User created successfully with ID: %s", id)

	Success(c, http.StatusCreated, "User created successfully", user)
}

// UpdateUser mengupdate user yang sudah ada
// Endpoint: PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	user, exists := users[id]
	if !exists {
		Error(c, http.StatusNotFound, "User not found", nil)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Validation error: %v", err)
		Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Update hanya field yang dikirim (partial update)
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.UpdatedAt = time.Now()

	users[id] = user
	h.logger.Infof("User updated successfully with ID: %s", id)

	Success(c, http.StatusOK, "User updated successfully", user)
}

// DeleteUser menghapus user
// Endpoint: DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if _, exists := users[id]; !exists {
		Error(c, http.StatusNotFound, "User not found", nil)
		return
	}

	delete(users, id)
	h.logger.Infof("User deleted successfully with ID: %s", id)

	Success(c, http.StatusOK, "User deleted successfully", nil)
}
