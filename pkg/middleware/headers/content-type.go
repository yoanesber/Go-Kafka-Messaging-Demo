package headers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

/**
 * ContentType is a middleware function that checks the Content-Type header of incoming requests.
 * It ensures that the Content-Type is set to `application/json` for POST, PUT, and PATCH requests.
 * If the Content-Type is not set correctly, it returns a 415 Unsupported Media Type error and aborts the request.
 * This middleware is useful for enforcing the expected content type for API requests.
 */
const (
	// contentTypeHeader is the header key for Content-Type
	contentTypeHeader = "Content-Type"
	// contentTypeJSON is the expected content type for JSON requests
	contentTypeJSON = "application/json"
)

func ContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		contentType := c.GetHeader(contentTypeHeader)

		// Only enforce for methods that require a body
		if method == http.MethodPost || method == http.MethodPut {
			if !strings.HasPrefix(contentType, contentTypeJSON) {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
					"error":   "Unsupported Media Type",
					"message": "Content-Type must be application/json",
				})
				return
			}
		}

		c.Next()
	}
}
