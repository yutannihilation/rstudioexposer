package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const (
	rStudioPubKeyURL = "http://localhost:8787/auth-public-key"
	rStudioLoginURL  = "http://localhost:8787/auth-do-sign-in"
	userAgent        = "Mozilla/5.0 (Windows NT 10.0; WOW64; rv:50.0) Gecko/20100101 Firefox/50.0"
)

func getPubkey() (*rsa.PublicKey, error) {
	req, err := http.NewRequest("GET", rStudioPubKeyURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to get pubkey from %s", rStudioPubKeyURL)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	pubkey, err := parsePubkey(string(body))
	if err != nil {
		return nil, err
	}
	return pubkey, nil
}

// parsePubkey parses the string returned by RStudio server
// The string is like this: 010001:DB1E3A8360F...
func parsePubkey(pubkey string) (*rsa.PublicKey, error) {
	s := strings.Split(pubkey, ":")
	if len(s) != 2 {
		return nil, fmt.Errorf("The format of the pubkey is wrong: %s", pubkey)
	}
	publicExponent, err := strconv.ParseInt(s[0], 16, 0)
	if err != nil {
		return nil, err
	}
	modulus := new(big.Int)
	_, success := modulus.SetString(s[1], 16)
	if !success {
		return nil, fmt.Errorf("Failed to parse modulus: %s", s[1])
	}
	return &rsa.PublicKey{N: modulus, E: int(publicExponent)}, nil
}

func getLoginSessionCookie(username, password string, pubkey *rsa.PublicKey) (*http.Cookie, error) {
	form := url.Values{}
	form.Add("persist", "0")
	form.Add("appUri", "")
	form.Add("clientPath", "")

	v, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey, []byte(username+"\n"+password))
	if err != nil {
		return nil, err
	}
	form.Add("v", string(base64.StdEncoding.EncodeToString(v)))

	req, err := http.NewRequest("POST", rStudioLoginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	localhostURL, _ := url.Parse("http://localhost")
	cookies := jar.Cookies(localhostURL)

	loginSessionCookie := cookies[0]
	if loginSessionCookie.Value == "" {
		return nil, fmt.Errorf("Failed to login!")
	}

	return loginSessionCookie, nil
}
