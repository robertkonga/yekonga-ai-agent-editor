package api

import (
	"log"
	"net/http"
	"time"
	"context"

	"github.com/gorilla/websocket"
	"github.com/yekonga/ai-agent/internal/agent"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for local development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents a connected WebSocket client
type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	engine *agent.Engine
}

func handleWebSocket(engine *agent.Engine, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		engine: engine,
	}

	log.Println("New WebSocket client connected")

	// Start go routines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump listens for messages from the frontend to the agent engine
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		log.Println("WebSocket client disconnected")
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		
		log.Printf("Received msg, starting Agent Loop for msg: %s", string(message))
		
		go func() {
			task := string(message)
			ctx := context.Background()
			
			// Start agent run loop
			err := c.engine.RunLoop(ctx, task, func(outputChunk string) {
				c.send <- []byte(outputChunk) // Stream agent thought / observations over WS
			})

			if err != nil {
				log.Printf("Agent Error: %v", err)
				c.send <- []byte("\n[Error processing task: " + err.Error() + "]")
			}
			c.send <- []byte("\n<done/>") // Sig done
		}()
	}
}

// writePump sends messages from the agent engine to the frontend
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second) // Ping interval
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
