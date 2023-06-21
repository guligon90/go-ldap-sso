package server

import (
	"log"
	"net/http"
	"text/template"

	"github.com/guligon90/go-ldap-sso/src/sso-srv/ldap"
	"github.com/guligon90/go-ldap-sso/src/sso-srv/sso"
)

var lsso sso.SSOer
var templates = template.New("")

const templatesDir = "src/sso-srv/templates/"

func tmplPath(title string) string {
	return templatesDir + title + ".html"
}

func init() {
	var err error
	lsso, err = ldap.NewLdapSSO()

	if err != nil {
		log.Fatalf("Error initializing ldap sso: %s", err)
	}

	templates = template.Must(
		template.ParseFiles(
			tmplPath("footer"),
			tmplPath("header"),
			tmplPath("login"),
		))
}

type TmplData struct {
	QueryString string
	Error       bool
}

func renderTemplate(w http.ResponseWriter, name string, p interface{}) {
	err := templates.ExecuteTemplate(w, name+".html", p)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
