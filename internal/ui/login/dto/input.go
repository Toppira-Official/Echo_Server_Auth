package dto

type LoginInput struct {
	Username string `json:"username" binding:"required,min=6,max=255" example:"echo_username"`
	Password string `json:"password" binding:"required,min=8,max=64" example:"echo_password"`
} //	@name	LoginInput
