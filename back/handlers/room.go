package handlers

import (
	"net/http"
	"quickstart/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomHandler struct {
    roomRepo models.RoomRepository
}

func NewRoomHandler(roomRepo models.RoomRepository) *RoomHandler {
    return &RoomHandler{roomRepo: roomRepo}
}

// CreateRoom godoc
// @Summary Create a new room
// @Schemes
// @Description Create a new chat room
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body models.Room true "Room object"
// @Success 201 {object} models.Room
// @Router /rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
    var room models.Room
    if err := c.ShouldBindJSON(&room); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.roomRepo.Create(&room); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, room)
}

// GetRooms godoc
// @Summary Get all rooms
// @Schemes
// @Description Get all chat rooms
// @Tags rooms
// @Accept json
// @Produce json
// @Success 200 {array} models.Room
// @Router /rooms [get]
func (h *RoomHandler) GetRooms(c *gin.Context) {
    rooms, err := h.roomRepo.FindAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, rooms)
}

// GetRoom godoc
// @Summary Get a room by ID
// @Schemes
// @Description Get a room by ID with users
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Success 200 {object} models.Room
// @Router /rooms/{id} [get]
func (h *RoomHandler) GetRoom(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
        return
    }
    
    room, err := h.roomRepo.FindByID(uint(id))
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, room)
}

// JoinRoom godoc
// @Summary Join a room
// @Schemes
// @Description Add user to a room
// @Tags rooms
// @Accept json
// @Produce json
// @Param roomId path int true "Room ID"
// @Param userId path int true "User ID"
// @Success 200 {object} models.Room
// @Router /rooms/{roomId}/join/{userId} [post]
func (h *RoomHandler) JoinRoom(c *gin.Context) {
    roomIdStr := c.Param("roomId")
    userIdStr := c.Param("userId")
    
    roomId, err := strconv.ParseUint(roomIdStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
        return
    }
    
    userId, err := strconv.ParseUint(userIdStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    
    if err := h.roomRepo.AddUser(uint(roomId), uint(userId)); err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Room or User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Load updated room with users
    room, err := h.roomRepo.FindByID(uint(roomId))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, room)
}