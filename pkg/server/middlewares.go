package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//VerifyHeader
func VerifyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		fmt.Println(token)
	}
}
