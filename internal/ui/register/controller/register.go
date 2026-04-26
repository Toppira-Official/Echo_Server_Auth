package controller

import (
	"auth/internal/application/usecase"
	"auth/internal/ui/register/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Register struct {
	registerUsecase *usecase.Register
}

func NewRegister(registerUsecase *usecase.Register) *Register {
	return &Register{registerUsecase: registerUsecase}
}

// Register godoc
//
//	@Summary		Register a new user account
//	@Description	Creates a new user using username and password, and returns access/refresh tokens.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.RegisterInput	true	"User registration input"
//	@Success		201		{object}	dto.RegisterOutput
//	@Failure		400		{object}	map[string]string	"Invalid input or business error"
//	@Router			/api/v1/auth/register [post]
func (r *Register) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	registerUsecaseInput := usecase.RegisterInput{
		Username:  input.Username,
		Password:  input.Password,
		UserAgent: c.Request.UserAgent(),
		IpAddress: c.ClientIP(),
	}
	registerUsecaseOutput, err := r.registerUsecase.Execute(c.Request.Context(), registerUsecaseInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterOutput{
		AccessToken:           registerUsecaseOutput.AccessToken,
		RefreshToken:          registerUsecaseOutput.RefreshToken,
		AccessTokenExpiresAt:  registerUsecaseOutput.AccessTokenExpiresAt.Format(time.RFC3339),
		RefreshTokenExpiresAt: registerUsecaseOutput.RefreshTokenExpiresAt.Format(time.RFC3339),
	})
}
