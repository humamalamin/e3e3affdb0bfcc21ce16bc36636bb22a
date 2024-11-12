package entity

import "github.com/golang-jwt/jwt/v5"

// Struct untuk value yang ada didalam token jwt
type JwtData struct {
	UserID float64 `json:"user_id"`
	jwt.RegisteredClaims
}
