package auth

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Iss string `json:"iss"`
	Exp int64  `json:"exp"`
	Jti string `json:"jti"`
	jwt.RegisteredClaims
}
