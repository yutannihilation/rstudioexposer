package main

import "net/http"

const (
	defaultUsername   = "rstudio"
	defaultPassword   = "rstudio"
	defaultListenPort = ":80"
)

func createRedirectHandler(cookie *http.Cookie) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "http://localhost:8787/", 302)
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
