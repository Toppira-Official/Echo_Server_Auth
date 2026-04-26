package middlewares

import (
	"net/http"

	"github.com/Ali127Dev/xerr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Error struct {
	logger *zap.Logger
}

func NewError(
	logger *zap.Logger,
) *Error {
	return &Error{logger: logger}
}

func RegisterErrorMiddleware(e *gin.Engine, m *Error) {
	e.Use(m.Handle())
}

func (m *Error) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		last := c.Errors.Last().Err
		if xe, ok := last.(*xerr.Error); ok {
			status := xe.Code().HTTPStatus()

			fields := []zap.Field{
				zap.String("code", string(xe.Code())),
				zap.Int("status", status),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			}

			if status >= 500 {
				m.logger.Error("server error", fields...)
				c.AbortWithStatusJSON(http.StatusInternalServerError, xerr.New(xerr.CodeInternalError))
				return
			}

			m.logger.Warn("client error", fields...)
			c.AbortWithStatusJSON(status, xe)
			return
		}

		m.logger.Error("unexpected error",
			zap.Error(last),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, xerr.New(xerr.CodeInternalError))
	}
}
