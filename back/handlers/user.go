package handlers

import (
	"net/http"
	"quickstart/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
    userRepo models.UserRepository
}

func NewUserHandler(userRepo models.UserRepository) *UserHandler {
    return &UserHandler{userRepo: userRepo}
}

// CreateUser godoc
// @Summary Create a new user
// @Schemes
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} models.User
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.userRepo.Create(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, user)
}

// GetUsers godoc
// @Summary Get all users
// @Schemes
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
    users, err := h.userRepo.FindAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Schemes
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    user, err := h.userRepo.FindByID(uint(id))
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, user)
}