package handlers

import (
    "net/http"
    "quickstart/models"
    
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    userRepo models.UserRepository
}

type LoginRequest struct {
    Name string `json:"name" binding:"required"`
}

type LoginResponse struct {
    Success bool         `json:"success"`
    Message string       `json:"message"`
    User    *models.User `json:"user,omitempty"`
}

func NewAuthHandler(userRepo models.UserRepository) *AuthHandler {
    return &AuthHandler{userRepo: userRepo}
}

// Login godoc
// @Summary Login user by name
// @Schemes
// @Description Login user using their name
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} LoginResponse
// @Failure 404 {object} LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, LoginResponse{
            Success: false,
            Message: "Invalid request: " + err.Error(),
        })
        return
    }

    // Find user by name
    users, err := h.userRepo.FindAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, LoginResponse{
            Success: false,
            Message: "Database error",
        })
        return
    }

    // Search for user with matching name
    var foundUser *models.User
    for _, user := range users {
        if user.Name == req.Name {
            foundUser = &user
            break
        }
    }

    if foundUser == nil {
        c.JSON(http.StatusNotFound, LoginResponse{
            Success: false,
            Message: "User not found",
        })
        return
    }

    c.JSON(http.StatusOK, LoginResponse{
        Success: true,
        Message: "Login successful",
        User:    foundUser,
    })
}