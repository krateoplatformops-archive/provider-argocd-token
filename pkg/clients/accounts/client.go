package accounts

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

// Login do a login with username and password credentials and returns the auth token.
func Login(opts *TokenProviderOptions, user, pass string) (string, error) {
	cli, err := NewTokenProvider(opts)
	if err != nil {
		return "", err
	}

	_, err = cli.CreateSession(user, pass)
	if err != nil {
		return "", err
	}

	return "", nil
}

// GenerateToken generate a token for the account with the specified name.
// expiresIn specify the duration before the token will expire; by default: no expiration.
func GenerateToken(opts *TokenProviderOptions, name string, expiresIn int64) (string, error) {
	cli, err := NewTokenProvider(opts)
	if err != nil {
		return "", err
	}

	res, err := cli.CreateTokenForAccount(name)
	if err != nil {
		return "", err
	}

	return res, nil
}

// TokenProviderOptions hold address, security, and other settings for the API client.
type TokenProviderOptions struct {
	ServerAddr string
	UserAgent  string
}

// TokenProvider defines an interface for interaction with an Argo CD server.
type TokenProvider interface {
	CreateSession(username, password string) (string, error)
	CreateTokenForAccount(name string) (string, error)
}

// NewTokenProvider creates a new ArgoCD token provider from a set of config options.
func NewTokenProvider(opts *TokenProviderOptions) (TokenProvider, error) {
	var res tokenProvider

	if opts.UserAgent != "" {
		res.UserAgent = opts.UserAgent
	}

	if opts.ServerAddr != "" {
		res.ServerAddr = opts.ServerAddr
	}
	// Make sure we got the server address and auth token from somewhere
	if res.ServerAddr == "" {
		return nil, errors.New("unspecified server address for Argo CD")
	}

	res.httpClient = &http.Client{}
	res.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &res, nil
}

type tokenProvider struct {
	ServerAddr string
	UserAgent  string
	httpClient *http.Client
}

func (tp tokenProvider) CreateSession(user, pass string) (string, error) {
	data := map[string]string{
		"username": user,
		"password": pass,
	}

	bin, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	var url string
	if strings.HasPrefix("http", tp.ServerAddr) {
		url = fmt.Sprintf("%s/api/v1/session", tp.ServerAddr)
	} else {
		url = fmt.Sprintf("https://%s/api/v1/session", tp.ServerAddr)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bin))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := tp.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	debug(httputil.DumpResponse(res, true))

	return "", nil
}

func (tp tokenProvider) CreateTokenForAccount(name string) (string, error) {
	return "", nil
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
