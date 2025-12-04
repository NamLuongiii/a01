package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Cho phép tất cả origins
	},
}

type WebSocketHandler struct{}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{}
}

func (wsh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Upgrade HTTP connection thành WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected!")

	// Gửi hello message
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello World! WebSocket server connected."))
	if err != nil {
		log.Println("Write error:", err)
		return
	}

	// Lắng nghe messages từ client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Printf("Received: %s", message)

		// Echo message trở lại client
		response := "Server received: " + string(message)
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}

	log.Println("Client disconnected!")
}