package types

type ListBook struct {
	Id           string `json:"id"`
	PaginationId int    `json:"pagination_id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	Description  string `json:"description"`
	IsBooked     bool   `json:"is_booked"`
	BookedUntil  string `json:"booked_until,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type CreateBook struct {
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateBook struct {
	Id          string `json:"id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Description string `json:"description" binding:"required"`
}
