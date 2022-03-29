package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"log"
	"net/http"
	"net/url"
)

var conf oauth2.Config

func checkedGetParam(query url.Values, key string) (string, bool) {
	param := query[key]
	if len(param) < 1 {
		return "", false
	}
	value := param[0]
	return value, true
}

func fillConfig() {

	id := flag.String("id", "", "the Github App ID to use for authentication")
	secret := flag.String("secret", "", "the Github App Secret to use for authentication")
	flag.Parse()
	if *id == "" || *secret == "" {
		log.Fatal("One one more of the following flags is missing: id, secret")
	}
	conf = oauth2.Config{
		ClientSecret: *secret,
		ClientID:     *id,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"repo"},
	}
}

func ensureParam(r *http.Request, param string) (string, error) {
	query := r.URL.Query()
	value, found := checkedGetParam(query, param)
	if !found {
		return "", fmt.Errorf("Query parameter '%v' missing", param)
	}

	return value, nil
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	state, err := ensureParam(r, "state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	redirectURL := conf.AuthCodeURL(state)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func makeIv(secret []byte, size int) []byte {
	if len(secret) >= size {
		return secret[:size]
	}
	doubled := append(secret, secret...)
	return makeIv(doubled, size)
}

func encrypt(data, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("could not create a cipher: %v", err)
	}
	encrypter := cipher.NewCFBEncrypter(block, makeIv(secret, block.BlockSize()))
	ciphered := make([]byte, len(data))
	encrypter.XORKeyStream(ciphered, data)
	return ciphered, nil
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code, err := ensureParam(r, "code")
	if err != nil {
		log.Panic(err)
	}

	state, err := ensureParam(r, "state")
	if err != nil {
		log.Panic(err)
	}

	token, err := conf.Exchange(r.Context(), code)
	if err != nil {
		log.Panicf("Could not exchange code for token: %v", err)
	}

	marshalled, err := json.Marshal(token)
	if err != nil {
		log.Panicf("Could not marshal token: %v", err)
	}

	encrypted, err := encrypt(marshalled, []byte(state))
	if err != nil {
		log.Panicf("Could not encrypt data: %v", err)
	}

	reader := bytes.NewReader(encrypted)

	_, err = http.Post("http://localhost:63231/token", "application/octet-stream", reader)
	if err != nil {
		log.Panicf("Could not post the token to the client server: %v", err)
	}
	_, err = w.Write([]byte(`<div style="display:flex;justify-content:center;align-items:center;font-size:3em;">Authentication finished. You can close this page.</div>`))
	if err != nil {
		log.Panicf("Could not write to the response: %v", err)
	}
}

func run() {
	http.HandleFunc("/auth", handleAuth)
	http.HandleFunc("/auth/callback", handleCallback)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server errored: %v\n", err)
	}
}

func main() {
	fillConfig()
	fmt.Println("Running server")
	run()
}
