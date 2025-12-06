package main

import (
	"log"
	docs "quickstart/docs"
	"quickstart/handlers"
	"quickstart/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @BasePath /api/v1

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

func main() {
  // Initialize Database
  var err error
  db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
  if err != nil {
    log.Fatal("Failed to connect to database:", err)
  }

  // Auto Migrate the schema
  db.AutoMigrate(&models.User{}, &models.Room{})

  // Initialize repositories
  userRepo := models.NewUserRepository(db)
  roomRepo := models.NewRoomRepository(db)

  // Initialize handlers
  userHandler := handlers.NewUserHandler(userRepo)
  roomHandler := handlers.NewRoomHandler(roomRepo)
  authHandler := handlers.NewAuthHandler(userRepo)
  wsHandler := handlers.NewWebSocketHandler()

  router := gin.Default()
  
  // Add CORS middleware
  router.Use(CORSMiddleware())
  
  docs.SwaggerInfo.BasePath = "/api/v1"

  v1 := router.Group("/api/v1")
  {
      // Auth routes
      auth := v1.Group("/auth")
      {
         auth.POST("/login", authHandler.Login)
      }

      // User routes
      users := v1.Group("/users")
      {
         users.POST("", userHandler.CreateUser)
         users.GET("", userHandler.GetUsers)
         users.GET("/:id", userHandler.GetUser)
      }
      
      // Room routes
      rooms := v1.Group("/rooms")
      {
         rooms.POST("", roomHandler.CreateRoom)
         rooms.GET("", roomHandler.GetRooms)
         rooms.GET("/:id", roomHandler.GetRoom)
         rooms.POST("/:roomId/join/:userId", roomHandler.JoinRoom)
      }
  }

  // WebSocket endpoint
  router.GET("/ws", wsHandler.HandleWebSocket)

  router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


  // Ping godoc
  // @Summary Ping the server
  // @Schemes
  // @Description Ping the server to check if it's alive
  // @Tags example
  // @Accept json
  // @Produce json
  // @Success 200 {object} map[string]string
  // @Router /ping [get]
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong🚀 hehe",
    })
  })

  router.Run() // listens on 0.0.0.0:8080 by default
  // http://localhost:8080/swagger/index.html#/example/get_example_helloworld
}