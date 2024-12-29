package types

type CreateBooking struct {
	BookId        string `json:"book_id" binding:"required"`
	CustomerName  string `json:"customer_name" binding:"required"`
	CustomerPhone string `json:"customer_phone" binding:"required"`
}

type GetBooking struct {
	Id            string `json:"id"`
	PaginationId  int    `json:"pagination_id"`
	BookId        string `json:"book_id"`
	BookTitle     string `json:"book_title"`
	BookAuthor    string `json:"book_author"`
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	BookedUntil   string `json:"booked_until"`
	IsReturned    bool   `json:"is_returned"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	UpdatedBy     string `json:"updated_by"`
	ReturnedAt    string `json:"returned_at,omitempty"`
}
