package main

import (
   "log"
   "github.com/gin-gonic/gin"
   docs "quickstart/docs"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "net/http"
   "gorm.io/gorm"
   "github.com/glebarez/sqlite"
)

// @BasePath /api/v1

// User model
type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
}

// Database instance
var db *gorm.DB

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

func main() {
  // Initialize Database
  var err error
  db, err = gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
  if err != nil {
    log.Fatal("Failed to connect to database:", err)
  }

  // Auto Migrate the schema
  db.AutoMigrate(&User{})

  router := gin.Default()
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