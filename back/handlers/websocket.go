package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Cho phép tất cả origins (development only)
    },
}

type Client struct {
    Conn   *websocket.Conn
    Send   chan []byte
    Hub    *Hub
    ID     string
    RoomID string
}

type Hub struct {
    Clients    map[*Client]bool
    Rooms      map[string]map[*Client]bool
    Broadcast  chan []byte
    Register   chan *Client
    Unregister chan *Client
}

type Message struct {
    Type    string `json:"type"`
    Content string `json:"content"`
    UserID  string `json:"user_id"`
    RoomID  string `json:"room_id"`
}

func NewHub() *Hub {
    return &Hub{
        Clients:    make(map[*Client]bool),
        Rooms:      make(map[string]map[*Client]bool),
        Broadcast:  make(chan []byte),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            h.Clients[client] = true
            log.Printf("Client %s connected", client.ID)
            
            // Thêm client vào room nếu có
            if client.RoomID != "" {
                if h.Rooms[client.RoomID] == nil {
                    h.Rooms[client.RoomID] = make(map[*Client]bool)
                }
                h.Rooms[client.RoomID][client] = true
                log.Printf("Client %s joined room %s", client.ID, client.RoomID)
            }
            
        case client := <-h.Unregister:
            if _, ok := h.Clients[client]; ok {
                delete(h.Clients, client)
                close(client.Send)
                
                // Xóa client khỏi room
                if client.RoomID != "" && h.Rooms[client.RoomID] != nil {
                    delete(h.Rooms[client.RoomID], client)
                    if len(h.Rooms[client.RoomID]) == 0 {
                        delete(h.Rooms, client.RoomID)
                    }
                }
                log.Printf("Client %s disconnected", client.ID)
            }
            
        case message := <-h.Broadcast:
            // Gửi message đến tất cả clients
            for client := range h.Clients {
                select {
                case client.Send <- message:
                default:
                    close(client.Send)
                    delete(h.Clients, client)
                }
            }
        }
    }
}

func (h *Hub) BroadcastToRoom(roomID string, message []byte) {
    if room, exists := h.Rooms[roomID]; exists {
        for client := range room {
            select {
            case client.Send <- message:
            default:
                close(client.Send)
                delete(h.Clients, client)
                delete(room, client)
            }
        }
    }
}

type WebSocketHandler struct {
    hub *Hub
}

func NewWebSocketHandler() *WebSocketHandler {
    hub := NewHub()
    go hub.Run()
    return &WebSocketHandler{hub: hub}
}

func (wsh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return
    }
    
    userID := c.Query("user_id")
    roomID := c.Query("room_id")
    
    client := &Client{
        Conn:   conn,
        Send:   make(chan []byte, 256),
        Hub:    wsh.hub,
        ID:     userID,
        RoomID: roomID,
    }
    
    client.Hub.Register <- client
    
    go client.writePump()
    go client.readPump()
}

func (c *Client) readPump() {
    defer func() {
        c.Hub.Unregister <- c
        c.Conn.Close()
    }()
    
    for {
        var msg Message
        err := c.Conn.ReadJSON(&msg)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        
        msg.UserID = c.ID
        
        // Xử lý các loại message khác nhau
        switch msg.Type {
        case "join_room":
            c.joinRoom(msg.RoomID)
        case "leave_room":
            c.leaveRoom()
        case "chat_message":
            c.broadcastToRoom(msg)
        }
    }
}

func (c *Client) writePump() {
    defer c.Conn.Close()
    
    for {
        select {
        case message, ok := <-c.Send:
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            c.Conn.WriteMessage(websocket.TextMessage, message)
        }
    }
}

func (c *Client) joinRoom(roomID string) {
    // Rời room cũ nếu có
    if c.RoomID != "" {
        c.leaveRoom()
    }
    
    // Join room mới
    c.RoomID = roomID
    if c.Hub.Rooms[roomID] == nil {
        c.Hub.Rooms[roomID] = make(map[*Client]bool)
    }
    c.Hub.Rooms[roomID][c] = true
    
    log.Printf("Client %s joined room %s", c.ID, roomID)
}

func (c *Client) leaveRoom() {
    if c.RoomID != "" && c.Hub.Rooms[c.RoomID] != nil {
        delete(c.Hub.Rooms[c.RoomID], c)
        if len(c.Hub.Rooms[c.RoomID]) == 0 {
            delete(c.Hub.Rooms, c.RoomID)
        }
        log.Printf("Client %s left room %s", c.ID, c.RoomID)
        c.RoomID = ""
    }
}

func (c *Client) broadcastToRoom(msg Message) {
    if c.RoomID != "" {
        messageBytes, _ := json.Marshal(msg)
        c.Hub.BroadcastToRoom(c.RoomID, messageBytes)
    }
}