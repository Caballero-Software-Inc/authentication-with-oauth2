package main

// tutorial: https://www.youtube.com/watch?v=OdyXIi6DGYw
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env")
		return
	}
	var (
		googleOauthConfig = &oauth2.Config{
			RedirectURL:  "http://localhost:8080/callback",
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		}
		//TODO: randomize
		randomState = "random"
	)

	handleHome := func(w http.ResponseWriter, r *http.Request) {
		var html = `<html> <body> <a href="/login">Google Log In</a> </body> </html>`
		fmt.Fprint(w, html)
	}

	handleLogin := func(w http.ResponseWriter, r *http.Request) {
		url := googleOauthConfig.AuthCodeURL(randomState)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

	handleCallback := func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != randomState {
			fmt.Println("state is not valid")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
		if err != nil {
			fmt.Printf("could not get token: %s\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		resp, err := http.Get("https://www.googleapis.com/oath2/v2/userinfo?access_token=" + token.AccessToken)
		if err != nil {
			fmt.Printf("could not create get request: %s\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("could not parse response: %s\n", err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		fmt.Fprintf(w, "Response: %s", content)
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil) //TODO get it from environement
}
