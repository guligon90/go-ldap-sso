package main

import (
	"log"
	"net/http"
	"os"

	"time"

	"github.com/guligon90/go-ldap-sso/src/sso-srv/ldap"
	"github.com/guligon90/go-ldap-sso/src/sso-srv/server"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	weblog "github.com/samitpal/goProbe/log"
)

func main() {
	const ssoSrvUrl = "localhost:8081"

	log.Println("Starting SSO server...")

	r := mux.NewRouter()

	var fh *os.File
	var err error

	wld := ldap.BaseConf.WeblogDir

	if wld != "" {
		fh, err = weblog.SetupWebLog(wld, time.Now())

		if err != nil {
			log.Fatalf("Failed to set up logging: %v", err)
		}
	} else {
		fh = os.Stdout // logs web accesses to stdout. May not be thread safe.
	}

	log.Println("Binding request handlers...")

	r.Handle("/sso", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(server.HandleSSOPostRequest))).Methods("POST")
	r.Handle("/sso", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(server.HandleSSOGetRequest))).Methods("GET")
	r.Handle("/logout", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(server.HandleLogoutRequest))).Methods("GET")
	r.Handle("/auth_token", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(server.HandleAuthTokenRequest))).Methods("POST")
	r.Handle("/test", handlers.CombinedLoggingHandler(fh, http.HandlerFunc(server.HandleTestRequest))).Methods("GET")

	http.Handle("/", r)

	log.Println("Done! Listening on https://" + ssoSrvUrl + "...")

	log.Fatal(http.ListenAndServeTLS(":8081", ldap.BaseConf.SSLCertPath, ldap.BaseConf.SSLKeyPath, nil))
	//log.Fatal(http.ListenAndServe(ssoSrvUrl, r))
}
