package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte // 待发送消息缓冲区
	UserID string      // 关联的用户 ID
}

// Hub 管理所有活跃的 WebSocket 连接
type Hub struct {
	// 所有已注册的客户端，key 为 userID
	Clients map[string]*Client
	mu      sync.RWMutex

	// 注册/注销/广播 channel
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func (h *Hub) NewHub() any {
	panic("unimplemented")
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte, 256),
	}
}

// Run 启动 Hub 事件循环（在独立 goroutine 中运行）
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:

					// 缓冲区满，关闭该连接
					close(client.Send)
					delete(h.Clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// 将收到的消息广播给所有人
		c.Hub.Broadcast <- message
	}
}

// WritePump 将待发送消息写入 WebSocket 连接
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Send
		if !ok {
			// Hub 关闭了 channel
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// SendToUser 向指定用户发送消息
func (h *Hub) SendToUser(userID string, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.mu.RLock()
	client, ok := h.Clients[userID]
	h.mu.RUnlock()

	if ok {
		select {
		case client.Send <- data:
		default:
			// 发送失败，连接可能已断开
		}
	}
	return nil
}
