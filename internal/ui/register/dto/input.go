package dto

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=6,max=255"`
	Password string `json:"password" binding:"required,min=8,max=64"`
} //	@name	RegisterInput
