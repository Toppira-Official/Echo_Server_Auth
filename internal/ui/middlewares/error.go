package middlewares

import (
	"net/http"

	"github.com/Ali127Dev/xerr"
	"github.com/gin-gonic/gin"
)

type Error struct{}

func NewError() *Error { return &Error{} }

func RegisterErrorMiddleware(e *gin.Engine, m *Error) {
	e.Use(m.Handle())
}

func (*Error) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		last := c.Errors.Last().Err
		if xe, ok := last.(*xerr.Error); ok {
			status := xe.Code().HTTPStatus()
			if status >= 500 {
				c.AbortWithStatusJSON(http.StatusInternalServerError, xerr.New(xerr.CodeInternalError))
				return
			}
			c.AbortWithStatusJSON(status, xe)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, xerr.New(xerr.CodeInternalError))
	}
}
