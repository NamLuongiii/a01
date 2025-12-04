package main

import (
	"log"
	"net/http"
	docs "quickstart/docs"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @BasePath /api/v1

// User model
type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
    Rooms []Room `json:"rooms" gorm:"many2many:user_rooms;"`
}

// Room model
type Room struct {
    ID          uint   `json:"id" gorm:"primaryKey"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Users       []User `json:"users" gorm:"many2many:user_rooms;"`
}

// Database instance
var db *gorm.DB

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context)  {
   g.JSON(http.StatusOK,"helloworld")
}

// CreateUser godoc
// @Summary Create a new user
// @Schemes
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User object"
// @Success 201 {object} User
// @Router /users [post]
func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := db.Create(&user).Error; err != nil {
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
// @Success 200 {array} User
// @Router /users [get]
func GetUsers(c *gin.Context) {
    var users []User
    if err := db.Find(&users).Error; err != nil {
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
// @Success 200 {object} User
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
    id := c.Param("id")
    var user User
    
    if err := db.First(&user, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// CreateRoom godoc
// @Summary Create a new room
// @Schemes
// @Description Create a new chat room
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body Room true "Room object"
// @Success 201 {object} Room
// @Router /rooms [post]
func CreateRoom(c *gin.Context) {
    var room Room
    if err := c.ShouldBindJSON(&room); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := db.Create(&room).Error; err != nil {
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
// @Success 200 {array} Room
// @Router /rooms [get]
func GetRooms(c *gin.Context) {
    var rooms []Room
    if err := db.Preload("Users").Find(&rooms).Error; err != nil {
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
// @Success 200 {object} Room
// @Router /rooms/{id} [get]
func GetRoom(c *gin.Context) {
    id := c.Param("id")
    var room Room
    
    if err := db.Preload("Users").First(&room, id).Error; err != nil {
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
// @Success 200 {object} Room
// @Router /rooms/{roomId}/join/{userId} [post]
func JoinRoom(c *gin.Context) {
    roomId := c.Param("roomId")
    userId := c.Param("userId")
    
    var room Room
    var user User
    
    if err := db.First(&room, roomId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
        return
    }
    
    if err := db.First(&user, userId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    if err := db.Model(&room).Association("Users").Append(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Load updated room with users
    db.Preload("Users").First(&room, roomId)
    c.JSON(http.StatusOK, room)
}

func main() {
  // Initialize Database
  var err error
  db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
  if err != nil {
    log.Fatal("Failed to connect to database:", err)
  }

  // Auto Migrate the schema
  db.AutoMigrate(&User{}, &Room{})

  router := gin.Default()
  
  // Add CORS middleware
  router.Use(CORSMiddleware())
  
  docs.SwaggerInfo.BasePath = "/api/v1"

  v1 := router.Group("/api/v1")
  {
      eg := v1.Group("/example")
      {
         eg.GET("/helloworld",Helloworld)
      }
      
      // User routes
      users := v1.Group("/users")
      {
         users.POST("", CreateUser)
         users.GET("", GetUsers)
         users.GET("/:id", GetUser)
      }
      
      // Room routes
      rooms := v1.Group("/rooms")
      {
         rooms.POST("", CreateRoom)
         rooms.GET("", GetRooms)
         rooms.GET("/:id", GetRoom)
         rooms.POST("/:roomId/join/:userId", JoinRoom)
      }
  }

  router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

  
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  router.Run() // listens on 0.0.0.0:8080 by default
  // http://localhost:8080/swagger/index.html#/example/get_example_helloworld
}