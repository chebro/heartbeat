package main

import (
	"github.com/go-chi/chi"
)

func router() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handleGetHome)
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/devices", handlePostDevices)
	})
	r.Route("/devices/{id}", func(r chi.Router) {
		r.Get("/", handleGetDevices)
	})
	return r
}
