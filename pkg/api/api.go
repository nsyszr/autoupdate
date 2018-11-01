package api

import "github.com/gorilla/mux"

// Service contains everything for running the auto update service
type API struct {
	r *mux.Router
}
