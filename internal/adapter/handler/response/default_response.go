package responseHandler

type DefaultSuccessResponse struct {
	Meta       Meta               `json:"meta"`
	Data       interface{}        `json:"data,omitempty"`
	Pagination PaginationResponse `json:"pagination,omitempty"`
}

type DefaultErrorResponse struct {
	Meta
}

type Meta struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type PaginationResponse struct {
	TotalCount int64 `json:"total_count"`
	PerPage    int   `json:"per_page"`
	Page       int   `json:"page"`
	TotalPages int   `json:"total_pages"`
}
