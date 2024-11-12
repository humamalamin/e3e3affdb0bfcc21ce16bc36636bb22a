package requestHandler

type ContentRequest struct {
	Title       string `json:"title" validate:"required"`
	Excerpt     string `json:"excerpt" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image       string `json:"image"`
	Tags        string `json:"tags"`
	CategoryID  int64  `json:"category_id" validate:"required"`
	Status      string `json:"status"`
}
