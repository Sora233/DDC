package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestid := c.GetHeader("X-Request-Id")
		if len(requestid) == 0 {
			requestid = uuid.New().String()
		}
		c.Set(KeyRequestID, requestid)
	}
}
