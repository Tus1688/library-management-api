package types

type CreateEmployee struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ListEmployee struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
