package updateserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insys-icom/autoupdate/pkg/repository"
)

// Handler contains everything for serving the update server
type Handler struct {
	prefix string
	r      *mux.Router
	dao    repository.DAO
}

// NewHandler creates a new update server handler
func NewHandler(prefix string, r *mux.Router, dao repository.DAO) *Handler {
	s := &Handler{
		prefix: prefix,
		r:      r,
		dao:    dao,
	}
	s.setupHandlerFuncs()
	return s
}

func (h *Handler) setupHandlerFuncs() {

	updateIndexHandler := http.HandlerFunc(h.handleUpdateIndex)
	updatePackageHandler := http.HandlerFunc(h.handleUpdatePackage)

	subRouter := h.r.PathPrefix(h.prefix).Subrouter()
	subRouter.Handle("/{tenant}/{target}/index.txt",
		h.clientAuthMiddleware(updateIndexHandler)).Methods("GET")
	subRouter.Handle("/{tenant}/{target}/packages/{package}",
		h.clientAuthMiddleware(updatePackageHandler)).Methods("GET")
}

func (h *Handler) clientAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// If no client certificate given, return 401
		if len(r.TLS.PeerCertificates) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		exists, err := h.dao.ExistsTargetBySlug(vars["tenant"], vars["target"])
		if err != nil {
			log.Printf("SQL error: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Fetch client certificate and check against tenant and target
		/* cert := r.TLS.PeerCertificates[0]
		if cert.Subject.CommonName != "client1" {
			w.WriteHeader(http.StatusForbidden)
			return
		} */

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "OK")
}

func (h *Handler) handleUpdateIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Deny potenial access to the PostgreSQL public schema
	if vars["tenant"] == "public" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Obtain a valid download list
	response, err := createDownloadList(h.dao, h.prefix, vars["tenant"], vars["target"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// HTTP Response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\n", response)
}

func (h *Handler) handleUpdatePackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "HandleUpdateIndex: %s\n  Tenant: %s\n  Target: %s\n  Package: %s\n",
		r.URL.Path, vars["tenant"], vars["target"], vars["package"])
}
