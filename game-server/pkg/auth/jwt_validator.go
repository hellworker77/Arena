package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTValidator struct {
	Cache     *JwksCache
	Authority string
	Audience  string
}

func (v *JWTValidator) ValidateToken(token string) (*JWTClaims, error) {
	parser := jwt.Parser{}
	unverifiedToken, _, err := parser.ParseUnverified(token, jwt.MapClaims{})

	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}

	kid, ok := unverifiedToken.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing kid")
	}

	pub, existing := v.Cache.GetKey(kid)
	if !existing {
		return nil, errors.New("unknown kid")
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return pub, nil
	}

	parsedToken, err := jwt.ParseWithClaims(
		token,
		&JWTClaims{},
		keyFunc,
		jwt.WithAudience(v.Audience),
		jwt.WithIssuer(v.Authority),
	)

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := parsedToken.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
