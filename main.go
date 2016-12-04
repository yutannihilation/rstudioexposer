package main

import "net/http"

import "strings"
import "net/url"

const (
	defaultUsername   = "rstudio"
	defaultPassword   = "rstudio"
	defaultListenPort = ":80"
)

func createRedirectHandler(cookie *http.Cookie) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, cookie)
		hostParts := strings.Split(r.Host, ":")
		publicURL := &url.URL{
			Scheme: "http",
			Host:   hostParts[0] + ":8787",
			Path:   "/",
		}

		http.Redirect(w, r, publicURL.String(), 302)
	}
}

func main() {
	pubkey, err := getPubkey()
	if err != nil {
		panic(err)
	}
	cookie, err := getLoginSessionCookie(defaultUsername, defaultPassword, pubkey)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", createRedirectHandler(cookie))

	http.ListenAndServe(defaultListenPort, nil)
}
