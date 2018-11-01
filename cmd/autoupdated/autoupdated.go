package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/insys-icom/autoupdate/pkg/repository/postgres"
	"github.com/insys-icom/autoupdate/pkg/updateserver"
	"github.com/urfave/negroni"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to PostgreSQL database
	connStr := "postgres://u4autoupdate:pw4autoupdate@localhost/autoupdate?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()

	// r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	// })
	/* r.HandleFunc("/", server.HandleIndex)
	r.HandleFunc("/autoupdate/v1/{tenant}/{target}/index.txt", server.HandleUpdateIndex).Methods("GET")
	r.HandleFunc("/autoupdate/v1/{tenant}/{target}/packages/{package}", server.HandleUpdatePackage) */

	// r.HandleFunc("/login", signin.HandleLogin)
	// r.HandleFunc("/consent", signin.HandleConsent)

	dao := postgres.NewDAO(db)
	updateserver.NewHandler("/autoupdate/v1", r, dao)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(r)

	// Setup HTTP server with TLS config
	caCert, err := ioutil.ReadFile("../../icomcloud-ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		// NoClientCert
		// RequestClientCert
		// RequireAnyClientCert
		// VerifyClientCertIfGiven
		// RequireAndVerifyClientCert
		ClientAuth: tls.VerifyClientCertIfGiven,
	}
	tlsConfig.BuildNameToCertificate()

	httpServer := &http.Server{
		Addr:      ":8443",
		Handler:   n,
		TLSConfig: tlsConfig,
	}

	httpServer.ListenAndServeTLS("../../localhost.crt", "../../localhost.key")
	// log.Fatal(http.ListenAndServe(":8080", n))
}
