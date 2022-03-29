package ghauth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func randomState() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 32)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}

func handleCallback(codeChan chan []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("could not read the callback body: %v", err)
		}

		codeChan <- body
	}
}

func runServer() chan []byte {
	ret := make(chan []byte)

	http.HandleFunc("/token", handleCallback(ret))
	go func() {
		err := http.ListenAndServe(":63231", nil)
		if err != nil {
			fmt.Printf("Server errored: %v\n", err)
		}
	}()

	return ret
}

func makeIv(secret []byte, size int) []byte {
	if len(secret) >= size {
		return secret[:size]
	}
	doubled := append(secret, secret...)
	return makeIv(doubled, size)
}

func decrypt(data, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("could not create a cipher: %v", err)
	}
	decrypter := cipher.NewCFBDecrypter(block, makeIv(secret, block.BlockSize()))
	decryptedData := make([]byte, len(data))
	decrypter.XORKeyStream(decryptedData, data)
	return decryptedData, nil
}

func GetToken() (string, error) {
	ret := runServer()
	state := randomState()
	url := fmt.Sprintf("http://localhost:8080/auth?state=%v", state)
	go func() {
		err := browser.OpenURL(url)
		if err != nil {
			log.Fatalf("could not open browser window: %v", err)
		}
	}()
	encryptedData := <-ret
	decryptedData, err := decrypt(encryptedData, []byte(state))
	if err != nil {
		return "", fmt.Errorf("could not decrypt received data: %v", err)
	}
	var token oauth2.Token
	err = json.Unmarshal(decryptedData, &token)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal the received data into a token: %v", err)
	}

	return token.AccessToken, nil
}
