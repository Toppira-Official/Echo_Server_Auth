package controller

import (
	"auth/internal/application/usecase"
	"auth/internal/ui/login/dto"
	"net/http"
	"time"

	_ "github.com/Ali127Dev/xerr"
	"github.com/gin-gonic/gin"
)

type Login struct {
	loginUsecase *usecase.Login
}

func NewLogin(loginUsecase *usecase.Login) *Login {
	return &Login{loginUsecase: loginUsecase}
}

// Login godoc
//
//	@Summary		Authenticate user
//	@Description	Authenticates a user using username and password and returns access and refresh tokens.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.LoginInput	true	"User login credentials"
//	@Success		200		{object}	dto.LoginOutput
//	@Failure		400		{object}	xerr.Error	"Invalid input"
//	@Failure		401		{object}	xerr.Error	"Invalid credentials"
//	@Failure		500		{object}	xerr.Error	"Internal server error"
//	@Router			/api/v1/auth/login [post]
func (l *Login) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.Error(err)
		return
	}

	loginUsecaseInput := usecase.LoginInput{
		Username:  input.Username,
		Password:  input.Password,
		UserAgent: c.Request.UserAgent(),
		IpAddress: c.ClientIP(),
	}
	loginUsecaseOutput, err := l.loginUsecase.Execute(c.Request.Context(), loginUsecaseInput)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, dto.LoginOutput{
		AccessToken:           loginUsecaseOutput.AccessToken,
		RefreshToken:          loginUsecaseOutput.RefreshToken,
		AccessTokenExpiresAt:  loginUsecaseOutput.AccessTokenExpiresAt.Format(time.RFC3339),
		RefreshTokenExpiresAt: loginUsecaseOutput.RefreshTokenExpiresAt.Format(time.RFC3339),
	})
}
