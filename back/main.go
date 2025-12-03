package main

import (
   "github.com/gin-gonic/gin"
   docs "quickstart/docs"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "net/http"
)

// @BasePath /api/v1

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

func main() {
  router := gin.Default()
  docs.SwaggerInfo.BasePath = "/api/v1"

  v1 := router.Group("/api/v1")
  {
      eg := v1.Group("/example")
      {
         eg.GET("/helloworld",Helloworld)
      }
  }

  router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

  
  router.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })

  router.Run() // listens on 0.0.0.0:8080 by default
}