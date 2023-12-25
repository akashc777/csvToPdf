package router

import (
	controllers "github.com/akashc777/csvToPdf/controllers/Templates"
	"net/http"

	customMiddleware "github.com/akashc777/csvToPdf/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(customMiddleware.RequestIDMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/login", controllers.Login)

	router.Get("/oauth_login/{oauth_provider}", controllers.OAuthLogin)
	router.Get("/oauth_callback/{oauth_provider}", controllers.OAuthCallback)

	router.Group(func(r chi.Router) {
		r.Use(customMiddleware.VerifyAuth)
		r.Post("/createTemplate", controllers.CreateTemplate)
		r.Get("/getTemplate", controllers.GetTemplateByName)
		r.Get("/getTemplateNames", controllers.GetTemplateNames)
		r.Post("/updateTemplate", controllers.UpdateTemplate)
		r.Delete("/deleteTemplate", controllers.DeleteTemplate)
	})

	return router
}
