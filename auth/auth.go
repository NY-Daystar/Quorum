package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"quorum/config"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GetClient read credentialsfile from google api
func GetClient(cfg *config.Config) (*http.Client, error) {
	b, err := os.ReadFile(cfg.CredentialsFile)
	if err != nil {
		return nil, err
	}

	conf, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/gmail.readonly")
	if err != nil {
		return nil, err
	}

	tok, err := getTokenFromFile(cfg.TokenFile)
	if err != nil {
		tok = getTokenFromWeb(conf, cfg.Port)
		saveToken(cfg.TokenFile, tok)
	}

	return conf.Client(context.Background(), tok), nil
}

// getTokenFromWeb listen on specific port
func getTokenFromWeb(config *oauth2.Config, port int64) *oauth2.Token {
	codeCh := make(chan string)

	// local server
	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(int(port)),
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintln(w, "Authorization successful! You can close this window.")
		codeCh <- code
	})

	go func() {
		_ = srv.ListenAndServe()
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Open this URL in your browser:\n", authURL)

	code := <-codeCh

	_ = srv.Close()

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		panic(err)
	}

	return tok
}

func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tok oauth2.Token
	err = json.NewDecoder(f).Decode(&tok)
	return &tok, err
}

func saveToken(path string, tok *oauth2.Token) {
	f, _ := os.Create(path)
	defer f.Close()
	json.NewEncoder(f).Encode(tok)
}
