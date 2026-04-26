package middlewares

import (
	"net/http"

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
		status := c.Writer.Status()
		if status < 400 {
			status = http.StatusInternalServerError
		}

		c.AbortWithStatusJSON(status, gin.H{
			"error": last.Error(),
		})
	}
}
