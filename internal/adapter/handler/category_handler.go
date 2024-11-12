package handler

import (
	requestHandler "latihan-portal-news/internal/adapter/handler/request"
	responseHandler "latihan-portal-news/internal/adapter/handler/response"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/internal/core/service"
	"latihan-portal-news/lib/conv"
	validatorLib "latihan-portal-news/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var defaultResponse responseHandler.DefaultSuccessResponse

type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error
	CreateCategory(c *fiber.Ctx) error
	EditCategoryByID(c *fiber.Ctx) error
	DeleteCategoryByID(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// DeleteCategoryByID implements CategoryHandler.
func (ch *categoryHandler) DeleteCategoryByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] DeleteCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParamString := c.Params("categoryID")
	id, err := conv.StringToInt64(idParamString)
	if err != nil {
		code = "[Handler] DeleteCategoryByID - 2"
		log.Errorw(code, err)

		errorResp.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(&errorResp)
	}

	err = ch.categoryService.DeleteCategoryByID(c.Context(), id)
	if err != nil {
		code = "[Handler] DeleteCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Data = nil
	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success delete category"

	return c.JSON(defaultResponse)
}

// EditCategoryByID implements CategoryHandler.
func (ch *categoryHandler) EditCategoryByID(c *fiber.Ctx) error {
	var req requestHandler.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] EditCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] EditCategoryByID - 2"
		log.Errorw(code, err)

		errorResp.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusUnprocessableEntity).JSON(&errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] EditCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnprocessableEntity).JSON(errorResp)
	}

	idParamString := c.Params("categoryID")
	id, err := conv.StringToInt64(idParamString)
	if err != nil {
		code = "[Handler] EditCategoryByID - 3"
		log.Errorw(code, err)

		errorResp.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(&errorResp)
	}
	reqEntity := entity.CategoryEntity{
		ID:    id,
		Title: req.Title,
		User:  entity.UserEntity{ID: int64(userID)},
	}
	err = ch.categoryService.EditCategoryByID(c.Context(), reqEntity)
	if err != nil {
		code = "[Handler] EditCategoryByID - 4"
		log.Errorw(code, err)

		errorResp.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(&errorResp)
	}
	defaultResponse.Data = nil
	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success edit category"

	return c.JSON(defaultResponse)
}

// CreateCategory implements CategoryHandler.
func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req requestHandler.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] CreateCategory - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] CreateCategory - 2"
		log.Errorw(code, err)

		errorResp.Status = false
		errorResp.Message = err.Error()

		return c.Status(fiber.StatusUnprocessableEntity).JSON(&errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] CreateCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnprocessableEntity).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err = ch.categoryService.CreateCategory(c.Context(), reqEntity)
	if err != nil {
		code = "[Handler] CreateCategory - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Data = nil
	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success"

	return c.Status(fiber.StatusCreated).JSON(defaultResponse)
}

// GetCategoryByID implements CategoryHandler.
func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] GetCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("categoryID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[Handler] GetCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}
	result, err := ch.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		code = "[Handler] GetCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponse := responseHandler.SuccessCategoryResponse{
		ID:            result.ID,
		Title:         result.Title,
		Slug:          result.Slug,
		CreatedByName: result.User.Name,
	}
	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success get category by ID"
	defaultResponse.Data = categoryResponse

	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// GetCategories implements CategoryHandler.
func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[Handler] GetCategories - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success get categories"

	categoryResponses := []responseHandler.SuccessCategoryResponse{}
	for _, cat := range results {
		categoryResponses = append(categoryResponses, responseHandler.SuccessCategoryResponse{
			ID:            cat.ID,
			Title:         cat.Title,
			Slug:          cat.Slug,
			CreatedByName: cat.User.Name,
		})
	}
	defaultResponse.Data = categoryResponses

	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{categoryService: categoryService}
}
