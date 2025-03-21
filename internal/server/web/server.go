package web

import (
	"log"
	"sender/internal/server/web/handlers"

	"github.com/gin-gonic/gin"
)

type WebServer struct {
	Port   string
	router *gin.Engine
}

// Create a new web server
func New(port string) WebServer {
	router := setupRoutes()
	return WebServer{
		Port:   port,
		router: router,
	}
}

func (ws *WebServer) Run() {
	// Start the server
	if err := ws.router.Run(":" + ws.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Print("Server running")
}

func (ws *WebServer) Router() *gin.Engine {
	return ws.router
}

// Function to set up routes
func setupRoutes() *gin.Engine {
	router := gin.Default()

	// Register routes
	router.GET("/health", handlers.HealthHandler)
	router.GET("/keys/generate", handlers.KeysGenerateHandler)

	return router
}
