package handler

import (
	"fmt"
	requestHandler "latihan-portal-news/internal/adapter/handler/request"
	responseHandler "latihan-portal-news/internal/adapter/handler/response"
	"latihan-portal-news/internal/core/domain/entity"
	"latihan-portal-news/internal/core/service"
	"latihan-portal-news/lib/conv"
	"latihan-portal-news/lib/pagination"
	validatorLib "latihan-portal-news/lib/validator"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	UpdateContent(c *fiber.Ctx) error
	DeleteContentByID(c *fiber.Ctx) error
	UploadImageCloudFlareR2(c *fiber.Ctx) error

	// Front
	GetContentByQuery(c *fiber.Ctx) error
	GetContentByCategoryID(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService   service.ContentService
	paginationHelper pagination.PaginationFunc
}

// UploadImageCloudFlareR2 implements ContentHandler.
func (cs *contentHandler) UploadImageCloudFlareR2(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] UploadImageCloudFlareR2 - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req requestHandler.FileUploadRequest
	file, err := c.FormFile("image")
	if err != nil {
		code = "[Handler] UploadImageCloudFlareR2 - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := c.SaveFile(file, fmt.Sprintf("./temp/content/%s", file.Filename)); err != nil {
		return err
	}
	req.Image = fmt.Sprintf("./temp/content/%s", file.Filename)

	reqEntity := entity.FileUploadEntity{
		Name: fmt.Sprintf("%d-%d", int64(userID), time.Now().UnixNano()),
		Path: req.Image,
	}

	imageUrl, err := cs.contentService.UploadImageCloudFlareR2(c.Context(), reqEntity)
	if err != nil {
		code = "[Handler] UploadImageCloudFlareR2 - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if req.Image != "" {
		err = os.Remove(req.Image)
		if err != nil {
			log.Errorw("[Handler] UploadImageCloudFlareR2 - 4", err)
		}
	}

	urlImageResp := map[string]interface{}{
		"urlImage": imageUrl,
	}
	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success upload image"
	defaultResponse.Data = urlImageResp

	return c.Status(fiber.StatusCreated).JSON(defaultResponse)
}

// GetContentByCategoryID implements ContentHandler.
func (cs *contentHandler) GetContentByCategoryID(c *fiber.Ctx) error {
	categoryIDParam := c.Params("categoryID")
	categoryID, err := conv.StringToInt64(categoryIDParam)
	if err != nil {
		code = "[Handler] GetContentByCategoryID - 1"
		log.Errorw(code, err)

		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	searchParams := c.Query("search")

	contents, err := cs.contentService.GetContentByCategoryID(c.Context(), categoryID, searchParams)
	if err != nil {
		code = "[Handler] GetContentByCategoryID - 2"
		log.Errorw(code, err)

		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success get contents"

	responseContent := []responseHandler.ContentResponse{}
	for _, val := range contents {
		responseContent = append(responseContent, responseHandler.ContentResponse{
			ID:           val.ID,
			Title:        val.Title,
			Excerpt:      val.Excerpt,
			Image:        val.Image,
			Status:       val.Status,
			CategoryName: val.Category.Title,
			Author:       val.User.Name,
			CreatedAt:    val.CreatedAt.Local().Format("02/Jan/2006"),
		})
	}

	defaultResponse.Data = responseContent

	return c.JSON(defaultResponse)
}

// GetContentByQuery implements ContentHandler.
func (cs *contentHandler) GetContentByQuery(c *fiber.Ctx) error {
	reqEntity := entity.QueryString{}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[Handler] GetContentByQuery - 1"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[Handler] GetContentByQuery - 2"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	order := "created_at"
	if c.Query("order") != "" {
		order = c.Query("order")
	}
	orderType := "desc"
	if c.Query("order_type") != "" {
		orderType = c.Query("order_type")
	}

	categoryIDString := c.Query("category_id")
	if categoryIDString != "" {
		categoryID, err := conv.StringToInt64(categoryIDString)
		if err != nil {
			code = "[Handler] GetContentByQuery - 3"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
		reqEntity.CategoryID = categoryID
	}

	reqEntity.Limit = limit
	reqEntity.Page = page
	reqEntity.Order = order
	reqEntity.OrderType = orderType

	results, err := cs.contentService.GetContents(c.Context(), &reqEntity)
	if err != nil {
		code = "[Handler] GetContentByQuery - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	pagin, err := cs.paginationHelper.AddPagination(len(results), page, limit)
	if err != nil {
		code = "[Handler] GetContentByQuery - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success get contents"

	responseContent := []responseHandler.ContentResponse{}
	for _, val := range results {
		responseContent = append(responseContent, responseHandler.ContentResponse{
			ID:           val.ID,
			Title:        val.Title,
			Excerpt:      val.Excerpt,
			Image:        val.Image,
			Status:       val.Status,
			CategoryName: val.Category.Title,
			Author:       val.User.Name,
			CreatedAt:    val.CreatedAt.Local().Format("02/Jan/2006"),
		})
	}

	defaultResponse.Data = responseContent[pagin.First:pagin.Last]
	defaultResponse.Pagination = responseHandler.PaginationResponse{
		TotalCount: int64(pagin.TotalCount),
		PerPage:    pagin.PerPage,
		Page:       pagin.Page,
		TotalPages: pagin.PageCount,
	}

	return c.JSON(defaultResponse)
}

// DeleteContentByID implements ContentHandler.
func (cs *contentHandler) DeleteContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] DeleteContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	contentIDParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(contentIDParam)
	if err != nil {
		code = "[Handler] DeleteContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = cs.contentService.DeleteContentByID(c.Context(), contentID)
	if err != nil {
		code = "[Handler] DeleteContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success delete content"
	defaultResponse.Data = nil

	return c.JSON(defaultResponse)
}

// UpdateContent implements ContentHandler.
func (cs *contentHandler) UpdateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] UpdateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req requestHandler.ContentRequest
	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] UpdateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] UpdateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnprocessableEntity).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Status:      req.Status,
		Tags:        tags,
	}

	contentIDParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(contentIDParam)
	if err != nil {
		code = "[Handler] UpdateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(&errorResp)
	}
	reqEntity.ID = contentID

	err = cs.contentService.UpdateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[Handler] UpdateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success update content"
	defaultResponse.Data = nil

	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// CreateContent implements ContentHandler.
func (cs *contentHandler) CreateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] CreateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req requestHandler.ContentRequest
	if err := c.BodyParser(&req); err != nil {
		code = "[Handler] CreateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnprocessableEntity).JSON(&errorResp)
	}

	if err = validatorLib.ValidateStruct(req); err != nil {
		code = "[Handler] CreateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnprocessableEntity).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Status:      req.Status,
		Tags:        tags,
	}

	err = cs.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[Handler] CreateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success create content"
	defaultResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultResponse)
}

// GetContentByID implements ContentHandler.
func (cs *contentHandler) GetContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] GetContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[Handler] GetContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := cs.contentService.GetContentByID(c.Context(), id)
	if err != nil {
		code = "[Handler] GetContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success get content by ID"

	resp := responseHandler.ContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
		CreatedAt:    result.CreatedAt.Local().Format("02/Jan/2006"),
		Status:       result.Status,
		CreatedByID:  result.CreatedByID,
		CategoryID:   result.CategoryID,
	}

	defaultResponse.Data = resp

	return c.Status(fiber.StatusOK).JSON(defaultResponse)
}

// GetContents implements ContentHandler.
func (cs *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[Handler] GetContents - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	reqEntity := entity.QueryString{
		Limit:      0,
		Page:       0,
		Order:      "",
		OrderType:  "",
		Search:     "",
		CategoryID: 0,
	}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[Handler] GetContentByQuery - 1"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[Handler] GetContentByQuery - 2"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	order := "created_at"
	if c.Query("order") != "" {
		order = c.Query("order")
	}
	orderType := "desc"
	if c.Query("order_type") != "" {
		orderType = c.Query("order_type")
	}

	categoryIDString := c.Query("category_id")
	if categoryIDString != "" {
		categoryID, err := conv.StringToInt64(categoryIDString)
		if err != nil {
			code = "[Handler] GetContentByQuery - 3"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
		reqEntity.CategoryID = categoryID
	}

	reqEntity.Limit = limit
	reqEntity.Page = page
	reqEntity.Order = order
	reqEntity.OrderType = orderType

	results, err := cs.contentService.GetContents(c.Context(), &reqEntity)
	if err != nil {
		code = "[Handler] GetContents - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultResponse.Meta.Status = true
	defaultResponse.Meta.Message = "success get contents"

	responseContent := []responseHandler.ContentResponse{}
	for _, val := range results {
		responseContent = append(responseContent, responseHandler.ContentResponse{
			ID:           val.ID,
			Title:        val.Title,
			Excerpt:      val.Excerpt,
			Image:        val.Image,
			Status:       val.Status,
			CategoryName: val.Category.Title,
			Author:       val.User.Name,
			CreatedAt:    val.CreatedAt.Local().Format("02/Jan/2006"),
		})
	}

	defaultResponse.Data = responseContent

	return c.JSON(defaultResponse)
}

func NewContentHandler(contentService service.ContentService, paginationHelper pagination.PaginationFunc) ContentHandler {
	return &contentHandler{contentService: contentService, paginationHelper: paginationHelper}
}
