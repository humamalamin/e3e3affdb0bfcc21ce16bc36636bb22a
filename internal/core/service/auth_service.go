package service

import (
	"context"
	"errors"
	"latihan-portal-news/config"
	"latihan-portal-news/internal/adapter/repository"
	"latihan-portal-news/internal/core/domain/entity"
	authLib "latihan-portal-news/lib/auth"
	"latihan-portal-news/lib/conv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var err error
var code string

type AuthService interface {
	GetUserByEmail(ctx context.Context, req entity.RequestLogin) (*entity.AccessTokenEntity, error)
}

type authService struct {
	authRepository repository.AuthRepository
	cfg            *config.Config
	jwtToken       authLib.Jwt
}

// GetUserByEmail implements AuthService.
func (a *authService) GetUserByEmail(ctx context.Context, req entity.RequestLogin) (*entity.AccessTokenEntity, error) {
	result, err := a.authRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code = "[Service] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if checkPass := conv.CheckPasswordHash(req.Password, result.Password); !checkPass {
		code = "[Service] GetUserByEmail - 2"
		err = errors.New("4001")
		log.Errorw(code, err)
		return nil, err
	}

	jwtData := entity.JwtData{
		UserID: float64(result.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			ID:        string(result.ID),
		},
	}

	accessToken, expiresAt, err := a.jwtToken.GenerateToken(&jwtData)
	if err != nil {
		code = "[Service] GetUserByEmail - 3"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.AccessTokenEntity{
		AccessToken: accessToken,
		ExpiredAt:   expiresAt,
	}, nil
}

func NewAuthService(
	authRepo repository.AuthRepository,
	cfg *config.Config,
	jwtToken authLib.Jwt,
) AuthService {
	return &authService{
		authRepository: authRepo,
		cfg:            cfg,
		jwtToken:       jwtToken,
	}
}
