package main

import "net/http"

const (
	defaultUsername   = "rstudio"
	defaultPassword   = "rstudio"
	defaultListenPort = ":80"
)

func createRedirectHandler(c string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", c)
		http.Redirect(w, r, "http://localhost:8787/", 302)
	}
}

func main() {
	pubkey, err := getPubkey()
	if err != nil {
		panic(err)
	}
	setCookieHeader, err := login(defaultUsername, defaultPassword, pubkey)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", createRedirectHandler(setCookieHeader))

	http.ListenAndServe(defaultListenPort, nil)
}
