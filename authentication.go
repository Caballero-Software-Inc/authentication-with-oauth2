package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html><body><a href="/login">Google Log In</a></body></html>`

	fmt.Fprintln(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var jsonContent map[string]interface{}

	json.Unmarshal([]byte(content), &jsonContent)
	var ok bool = jsonContent["verified_email"] == true
	fmt.Fprintf(w, "Content: %s\n", content)

	if ok {
		uploadUser(jsonContent["id"].(string), jsonContent["email"].(string))
	} else {
		fmt.Println("Authentication failed because the email was not verified.")
	}
}
