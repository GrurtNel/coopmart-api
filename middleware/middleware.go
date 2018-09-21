package middleware

import (
	"feedback/x/rest"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var errResponse = map[string]interface{}{
					"error":  err.(error).Error(),
					"status": "error",
				}
				if httpError, ok := err.(rest.IHttpError); ok {
					c.JSON(httpError.StatusCode(), errResponse)
				} else {
					c.JSON(500, errResponse)
				}
			}
		}()
		c.Next()
	}
}

func AddHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		//remember
		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(200)
			return
		}
		c.Next()
	}
}
