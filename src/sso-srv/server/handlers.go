package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/guligon90/go-ldap-sso/src/sso-srv/sso"
)

// HandleSSOGetRequest presents the login form
func HandleSSOGetRequest(w http.ResponseWriter, r *http.Request) {
	err := false

	if r.URL.Query().Get("auth_error") != "" {
		err = true
	}

	log.Println(r.URL.Query().Get("s_url"))

	tmplData := TmplData{QueryString: r.URL.Query().Get("s_url"), Error: err}

	renderTemplate(w, "login", &tmplData)
}

// HandleSSOPostRequest sets the sso cookie.
func HandleSSOPostRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	p_uri := r.PostFormValue("query_string")

	u, g, err := lsso.Auth(r.PostFormValue("username"), r.PostFormValue("password"))

	if u != nil {
		vh := lsso.CTValidHours()
		exp := time.Now().Add(time.Hour * time.Duration(vh)).UTC()
		tok, _ := lsso.BuildJWTToken(*u, *g, exp)
		c := lsso.BuildCookie(tok, exp)

		http.SetCookie(w, &c)
		http.Redirect(w, r, p_uri, 301)

		return
	}

	if err != nil {
		if sso.Err401Map[err] {
			log.Println(err)
			http.Redirect(w, r, fmt.Sprintf("/sso?s_url=%s&auth_error=true", p_uri), 301)

			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Not able to service this request. Please try again later.")

		return

	}
}

// HandleAuthTokenRequest generates the raw jwt token and sends it across.
func HandleAuthTokenRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u, g, err := lsso.Auth(r.PostFormValue("username"), r.PostFormValue("password"))

	if u != nil {
		tok, _ := lsso.BuildJWTToken(*u, *g, time.Now().Add(time.Hour*time.Duration(lsso.CTValidHours())).UTC())
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tok)

		return
	}
	if err != nil {
		if sso.Err401Map[err] {
			log.Println(err)
			fmt.Fprintf(w, "Unauthorized.")

			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Not able to service the request. Please try again later.")

		return

	}
}

// HandleLogoutRequest function invalidates the sso cookie.
func HandleLogoutRequest(w http.ResponseWriter, r *http.Request) {
	expT := time.Now().Add(time.Hour * time.Duration(-1))
	lc := lsso.Logout(expT)

	http.SetCookie(w, &lc)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "You have been logged out.")

	return
}

// HandleTestRequest function is just for the purpose of testing.
func HandleTestRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "You have visited a test page.")

	return
}
