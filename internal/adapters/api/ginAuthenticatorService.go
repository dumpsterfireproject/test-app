package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinAuthenicatorService interface {
	AuthenticateToken(c *gin.Context)
}

type GinAuthenicatorServiceImpl struct {
}

func NewGinAuthenicatorService() GinAuthenicatorService {
	return &GinAuthenicatorServiceImpl{}
}

// For demo, right now we accept any token
func (auth *GinAuthenicatorServiceImpl) AuthenticateToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
