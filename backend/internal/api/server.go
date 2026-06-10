package api

import (
	"log"
	"net/http"
	"github.com/yekonga/ai-agent/internal/agent"
)

type Server struct {
	addr   string
	engine *agent.Engine
}

func NewServer(addr string, engine *agent.Engine) *Server {
	return &Server{
		addr:   addr,
		engine: engine,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Setting up API routes
	s.registerRoutes(mux)

	log.Printf("Server listening on %s", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	// Enable CORS for frontend
	
	mux.HandleFunc("/api/health", s.handleHealth)
	
	// WebSocket endpoint for the agent chat/IDE
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(s.engine, w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Simple CORS handling
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
