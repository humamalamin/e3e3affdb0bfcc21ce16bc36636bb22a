package handler

import (
	requestHandler "latihan-portal-news/internal/adapter/handler/request"
	responseHandler "latihan-portal-news/internal/adapter/handler/response"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/internal/core/service"
	validatorLib "latihan-portal-news/lib/validator"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var err error
var code string
var errorResp responseHandler.DefaultErrorResponse
var validate = validator.New()

type AuthHandler interface {
	Login(c *fiber.Ctx) error
}

type authHandler struct {
	authService service.AuthService
}

// Login implements AuthHandler.
func (a *authHandler) Login(c *fiber.Ctx) error {
	req := requestHandler.LoginRequest{}
	resp := responseHandler.SuccessAuthResponse{}

	if err = c.BodyParser(&req); err != nil {
		code = "[Handler] Login - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] Login - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnprocessableEntity).JSON(errorResp)
	}

	reqLogin := entity.RequestLogin{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := a.authService.GetUserByEmail(c.Context(), reqLogin)
	if err != nil {
		if err.Error() == "4001" {
			code = "[Handler] LoginByEmail - 3"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Email / Password invalid"

			return c.Status(fiber.StatusUnauthorized).JSON(&errorResp)
		}
		code = "[Handler] Login - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	resp.Meta.Status = true
	resp.Meta.Message = "success login"
	resp.AccessToken = result.AccessToken
	resp.ExpiredAt = result.ExpiredAt

	return c.JSON(resp)
}

func NewAuthhandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}
