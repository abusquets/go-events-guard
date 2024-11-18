// file: internal/infrastructure/server/mux_server.go

package server

import (
	"context"
	"eventsguard/internal/infrastructure/config"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"go.uber.org/fx"
)

type webServer struct {
	config *config.AppConfig
	mux    *http.ServeMux
	server *http.Server
}

// NewMuxServer inicialitza un nou servidor amb el mux.
func NewMuxServer(cfg *config.AppConfig) WebServer {
	mux := http.NewServeMux()

	// Afegim una ruta de salut
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Healthy")
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Configuració del servidor HTTP
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Ajusta això a l'URL del teu frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	address := fmt.Sprintf(":%d", cfg.ServerPort)

	// Cream el servidor HTTP amb un handler per CORS
	httpServer := &http.Server{
		Addr:    address,
		Handler: c.Handler(mux),
	}

	return &webServer{
		config: cfg,
		mux:    mux,
		server: httpServer,
	}
}

// Start inicia el servidor HTTP.
func (s *webServer) Start() {
	log.Printf("Starting server on %s\n", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", s.server.Addr, err)
	}
}

// Shutdown tanca el servidor HTTP.
func (s *webServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down the server gracefully...")
	return s.server.Shutdown(ctx)
}

// GetMux retorna el ServeMux del servidor.
func (s *webServer) GetMux() *http.ServeMux {
	return s.mux
}

func (s *webServer) RegisterHooks(lifecycle fx.Lifecycle) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go s.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Shutting down ApiWebServer...")
				return s.Shutdown(ctx)
			},
		},
	)
}
