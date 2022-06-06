package main

import (
	"github.com/go-chi/chi"
)

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handleGetHome)
		r.Get("/edit", handleGetEdit)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/devices", handlePostDevices)
		r.Delete("/devices/{id}", handleDeleteDevice)
	})
	r.Route("/devices/{id}", func(r chi.Router) {
		r.Get("/", handleGetDevice)
	})
	return r
}
