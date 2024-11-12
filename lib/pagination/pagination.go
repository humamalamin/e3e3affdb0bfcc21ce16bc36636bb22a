package pagination

import (
	"latihan-portal-news/internal/core/domain/entity"
	"math"
)

type PaginationFunc interface {
	AddPagination(totalData int, page, perPage int) (*entity.Page, error)
}

type Options struct{}

// AddPagination implements Pagination.
func (o *Options) AddPagination(totalData int, page, perPage int) (*entity.Page, error) {

	newPage := page

	if newPage <= 0 {
		return nil, ErrorPage
	}

	limitData := 10
	if perPage != 0 {
		limitData = perPage
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(limitData)))

	last := (newPage * limitData)
	first := last - limitData

	if totalData < last {
		last = totalData
	}

	zeroPage := &entity.Page{PageCount: 1, Page: newPage}
	if totalPage == 0 && newPage == 1 {
		return zeroPage, nil
	}

	if newPage > totalPage {
		return nil, ErrorMaxPage
	}

	pages := &entity.Page{
		PageCount:  totalPage,
		Page:       newPage,
		TotalCount: totalData,
		PerPage:    limitData,
		First:      first,
		Last:       last,
	}

	return pages, nil
}

func NewPagination() PaginationFunc {
	pagination := new(Options)

	return pagination
}
