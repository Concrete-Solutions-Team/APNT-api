package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/slupx/smartest-backend/internal/auth"
)

func (s *Server) MountEndpoints() {
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://smartest-ui.vercel.app/"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	s.router.Route("/auth", func(r chi.Router) {
		r.Post("/register", s.authHandler.Register)
		r.Post("/login", s.authHandler.Login)
		r.Post("/logout", s.authHandler.Logout)
		r.With(s.authHandler.Identify).Get("/me", s.authHandler.Me)
	})

	s.router.Group(func(r chi.Router) {
		r.Use(s.authHandler.Identify)

		r.Route("/tests", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(s.authHandler.RequireRole(auth.RoleTeacher))
				r.Get("/", s.testHandler.ListTests)
				r.Post("/", s.testHandler.CreateTest)
			})
			r.Get("/join", s.testHandler.JoinTest)
			r.Post("/submit", s.testHandler.SubmitTest)
			r.Get("/results", s.testHandler.GetResults)
		})
	})

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("aaa?\n"))
	})
}
