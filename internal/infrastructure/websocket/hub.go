package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSEvent representa un evento que se envía a los clientes WebSocket
type WSEvent struct {
	Type    string      `json:"type"`    // new_alert, alert_updated, battery_update, device_status
	Payload interface{} `json:"payload"` // Datos específicos del evento
}

// Client representa una conexión WebSocket activa
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// Hub mantiene el conjunto de clientes activos y broadcastea mensajes
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var (
	// GlobalHub es la instancia global del Hub de WebSocket
	GlobalHub *Hub

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Permitir todas las conexiones (desarrollo)
		},
	}
)

// NewHub crea una nueva instancia del Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// InitGlobalHub inicializa el Hub global y lo pone a correr
func InitGlobalHub() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
	log.Println("WebSocket Hub iniciado correctamente")
}

// Run ejecuta el loop principal del Hub para procesar registros y broadcasts
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket: Cliente conectado (total: %d)", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket: Cliente desconectado (total: %d)", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Si el canal está lleno, desconectamos al cliente lento
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastEvent envía un evento a todos los clientes conectados
func (h *Hub) BroadcastEvent(event WSEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error serializando evento WebSocket: %v", err)
		return
	}
	h.broadcast <- data
}

// Broadcast es un helper global para enviar eventos desde cualquier parte de la app
func Broadcast(eventType string, payload interface{}) {
	if GlobalHub == nil {
		return
	}
	GlobalHub.BroadcastEvent(WSEvent{
		Type:    eventType,
		Payload: payload,
	})
}

// HandleWebSocket maneja las conexiones WebSocket HTTP
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al hacer upgrade de WebSocket: %v", err)
		return
	}

	client := &Client{
		hub:  GlobalHub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	GlobalHub.register <- client

	// Goroutine para escribir mensajes al cliente
	go client.writePump()
	// Goroutine para leer mensajes del cliente (mantener conexión viva)
	go client.readPump()
}

// writePump envía mensajes desde el hub al cliente WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second) // Ping cada 30s para mantener viva la conexión
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump lee mensajes del cliente (principalmente para detectar desconexiones)
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error de lectura: %v", err)
			}
			break
		}
	}
}
