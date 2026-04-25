package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (s *Server) MountEndpoints() {
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.router.Use(s.authHandler.Identify)

	s.router.Route("/auth", func(r chi.Router) {
		r.Post("/register", s.authHandler.Register)
		r.Post("/login", s.authHandler.Login)
		r.Post("/logout", s.authHandler.Logout)
		r.Get("/me", s.authHandler.Me)
	})

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("aaa?\n"))
	})

	// s.router.Post("/subscribe", s.subHandler.Subscribe)
	// s.router.Get("/confirm/{token}", s.subHandler.Confirm)
	// s.router.Get("/unsubscribe/{token}", s.subHandler.Unsubscribe)
	// s.router.Get("/subscriptions", s.subHandler.Subscriptions)
}
