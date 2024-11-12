package authLib

import (
	"fmt"
	"latihan-portal-news/config"
	"latihan-portal-news/internal/core/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	GenerateToken(data *entity.JwtData) (string, int64, error)
	VerifyAccessToken(token string) (*entity.JwtData, error)
}

type Options struct {
	signingKey          string
	issuer              string
	accessTokenDuration int64
}

// GenerateToken implements Jwt.
func (o *Options) GenerateToken(data *entity.JwtData) (string, int64, error) {
	now := time.Now().Local()
	expiresAt := now.Add(time.Hour * 1)
	data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	data.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	acToken := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	accessToken, err := acToken.SignedString([]byte(o.signingKey))
	if err != nil {
		return "", 0, err
	}

	return accessToken, int64(data.RegisteredClaims.ExpiresAt.Unix()), nil
}

// VerifyAccessToken implements Jwt.
func (o *Options) VerifyAccessToken(token string) (*entity.JwtData, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		}

		return []byte(o.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			return nil, err
		}

		jwtData := &entity.JwtData{
			UserID: claims["user_id"].(float64),
		}

		return jwtData, nil
	} else {
		return nil, fmt.Errorf("Token is invalid")
	}
}

func NewJwt(cfg *config.Config) Jwt {
	opt := new(Options)
	opt.issuer = cfg.App.JwtIssuer
	opt.signingKey = cfg.App.JwtSecretKey
	opt.accessTokenDuration = int64(cfg.App.JwtDurationAccessKey)

	return opt
}
