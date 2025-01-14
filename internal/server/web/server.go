package web

import (
	"log"
	"net/http"
	"sender/internal/server/web/handlers"

	"github.com/gorilla/mux"
)

type WebServer struct {
	Port   string
	routes *mux.Router
}

// Создание нового веб-сервера
func New(port string) WebServer {
	routes := identifyRoutes()
	return WebServer{
		Port:   port,
		routes: routes,
	}
}

func (ws *WebServer) Run() {
	// Запуск сервера
	if err := http.ListenAndServe(":"+ws.Port, ws.routes); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Print("Server run")
}

// Функция для настройки маршрутов
func identifyRoutes() *mux.Router {
	router := mux.NewRouter()

	// Регистрация маршрутов
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	router.HandleFunc("/keys/generate", handlers.KeysGenerateHandler).Methods("GET")

	return router
}
