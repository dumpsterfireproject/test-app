package api

import (
	"fmt"
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

// For demo, right now we accept any non-empty token
func (auth *GinAuthenicatorServiceImpl) AuthenticateToken(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	fmt.Printf("HEADER (%s)", header)
	if !strings.HasPrefix(header, "Bearer ") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token := header[7:]
	if len(token) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
