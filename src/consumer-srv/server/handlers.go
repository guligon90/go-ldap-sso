package server

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/guligon90/go-ldap-sso/src/util"
)

var parsedPubKey *rsa.PublicKey

func init() {
	key, _ := ioutil.ReadFile("../key_pair/demo.rsa.pub") // this is the public key of the login (github.com/guligon90/go-ldap-sso) server
	parsedPubKey, _ = jwt.ParseRSAPublicKeyFromPEM(key)
}

func HandleCookieCheck(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("SSO_C")

	if err == http.ErrNoCookie {
		// we redirect to the login service setting the appropriate s_url to come back after auth.
		http.Redirect(w, r, "https://127.0.0.1:8081/sso?s_url=https://127.0.0.1:8082/cookie", 301)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parts := strings.Split(strings.Split(c.String(), "=")[1], ".")
	err = jwt.SigningMethodRS512.Verify(strings.Join(parts[0:2], "."), parts[2], parsedPubKey)

	if err != nil {
		log.Fatalf("[%v] Error while verifying key: %v", strings.Split(c.String(), "=")[1], err)
	}

	tokenString := strings.Split(c.String(), "=")[1]
	token, err := jwt.ParseWithClaims(tokenString, &util.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsedPubKey, nil
	})

	claims, ok := token.Claims.(*util.CustomClaims) // claims.User and claims.Roles are what we are interested in.

	if ok && token.Valid {
		fmt.Printf("User: %v Roles: %v Tok_Expires: %v \n", claims.User, claims.Roles, claims.StandardClaims.ExpiresAt)

	} else {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "You have visited a cookietest page.\n\n")
	fmt.Fprintf(w, "User: %v, Roles: %v, Tok_Expires: %v\n", claims.User, claims.Roles, claims.StandardClaims.ExpiresAt)

	return
}

func HandleAuthTokCheck(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("Authorization")

	if h == "" {
		http.Error(w, "No Authorization header", http.StatusInternalServerError)
		return
	}

	parts := strings.Split(strings.Split(h, " ")[1], ".")
	err := jwt.SigningMethodRS512.Verify(strings.Join(parts[0:2], "."), parts[2], parsedPubKey)

	if err != nil {
		log.Fatalf("[%v] Error while verifying key: %v", strings.Split(h, "=")[1], err)
	}

	tokenString := strings.Split(h, " ")[1]
	token, err := jwt.ParseWithClaims(tokenString, &util.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsedPubKey, nil
	})

	claims, ok := token.Claims.(*util.CustomClaims) // claims.User and claims.Roles are what we are interested in.

	if ok && token.Valid {
		fmt.Printf("User: %v Roles: %v Tok_Expires: %v \n", claims.User, claims.Roles, claims.StandardClaims.ExpiresAt)

	} else {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "You have visited a auth tokentest page.\n\n")
	fmt.Fprintf(w, "User: %v, Roles: %v, Tok_Expires: %v\n", claims.User, claims.Roles, claims.StandardClaims.ExpiresAt)

	return
}
