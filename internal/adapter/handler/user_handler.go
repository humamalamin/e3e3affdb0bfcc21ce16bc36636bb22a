package handler

import (
	requestHandler "latihan-portal-news/internal/adapter/handler/request"
	responseHandler "latihan-portal-news/internal/adapter/handler/response"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/internal/core/service"
	validatorLib "latihan-portal-news/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserHandler interface {
	UpdatePassword(c *fiber.Ctx) error
	GetUserByID(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

// GetUserByID implements UserHandler.
func (u *userHandler) GetUserByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] GetUserByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	user, err := u.userService.GetUserByID(c.Context(), int64(userID))
	if err != nil {
		code = "[Handler] GetUserByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success get user by ID"

	resp := responseHandler.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	defaultResponse.Data = resp

	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// UpdatePassword implements UserHandler.
func (u *userHandler) UpdatePassword(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] UpdatePassword - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req requestHandler.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnprocessableEntity).JSON(errorResp)
	}

	err = u.userService.UpdatePassword(c.Context(), req.NewPassword, int64(userID))
	if err != nil {
		code = "[Handler] UpdatePassword - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success update password"
	defaultResponse.Data = nil
	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}
