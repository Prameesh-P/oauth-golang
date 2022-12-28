package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/return",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.goolgeapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	RandomState = "RandomString"
)

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/return", ReturnHandler)
	http.ListenAndServe(":8080", nil)
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	html := `<html><body><a href="/login">Google Log In</a></body></html>`
	fmt.Fprintf(w, html)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(RandomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func ReturnHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != RandomState {
		fmt.Println("state not valid..")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Printf("Could NOT get token %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("Could NOT Create request: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Could NOT parse response: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "Response: %s", content)
}
