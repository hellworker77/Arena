package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JwkKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

type Jwks struct {
	Keys []JwkKey `json:"keys"`
}

func FetchJwks(authority string) (*Jwks, error) {
	url := fmt.Sprintf("%s/.well-known/auth.json", authority)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var jwks Jwks
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}
