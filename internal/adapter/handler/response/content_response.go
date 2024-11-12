package responseHandler

type ContentResponse struct {
	ID           int64    `json:"id"`
	Title        string   `json:"title"`
	Excerpt      string   `json:"excerpt"`
	Description  string   `json:"description,omitempty"`
	Image        string   `json:"image"`
	Tags         []string `json:"tags,omitempty"`
	CategoryName string   `json:"category_name"`
	Author       string   `json:"author"`
	CreatedAt    string   `json:"created_at"`
	Status       string   `json:"status"`
	CreatedByID  int64    `json:"created_by_id,omitempty"`
	CategoryID   int64    `json:"category_id,omitempty"`
}
