package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/guligon90/go-ldap-sso/src/consumer-srv/server"

	"github.com/gorilla/mux"
)

func main() {
	const consumerUrl = "localhost:8082"

	log.Println("Starting app server...")

	r := mux.NewRouter()

	log.Println("Binding request handlers...")

	r.HandleFunc("/cookie", server.HandleCookieCheck)
	r.HandleFunc("/auth_token", server.HandleAuthTokCheck)

	http.Handle("/", r)

	certFilePath, _ := filepath.Abs("./ssl_certs/cert.pem")
	certKeyPath, _ := filepath.Abs("./ssl_certs/key.pem")

	log.Println("Done! Listening to https://" + consumerUrl + "...")
	//log.Fatal(http.ListenAndServe(consumerUrl, r))

	log.Fatal(http.ListenAndServeTLS(":8082", certFilePath, certKeyPath, nil))
}
