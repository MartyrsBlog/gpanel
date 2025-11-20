package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			
			// 根据错误类型返回不同的状态码
			var statusCode int
			var message string

			switch err.(type) {
			case *ValidationError:
				statusCode = http.StatusBadRequest
				message = err.Error()
			case *AuthError:
				statusCode = http.StatusUnauthorized
				message = err.Error()
			case *ForbiddenError:
				statusCode = http.StatusForbidden
				message = err.Error()
			case *NotFoundError:
				statusCode = http.StatusNotFound
				message = err.Error()
			default:
				statusCode = http.StatusInternalServerError
				message = "Internal server error"
			}

			c.JSON(statusCode, gin.H{
				"code":    statusCode,
				"message": message,
				"data":    nil,
			})
		}
	}
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// AuthError 认证错误
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// ForbiddenError 权限错误
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

// NotFoundError 资源不存在错误
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}