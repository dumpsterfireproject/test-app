package api

import (
	"github.com/gin-gonic/gin"
)

type GinAuthenticatorEndpoint struct {
	service GinAuthenicatorService
}

func NewGinAuthenticatorEndpoint() *GinAuthenticatorEndpoint {
	service := NewGinAuthenicatorService()
	return &GinAuthenticatorEndpoint{service: service}
}

func (auth *GinAuthenticatorEndpoint) AuthenticatedRouterGroup(router *gin.Engine, path string) *gin.RouterGroup {
	return router.Group(path, auth.service.AuthenticateToken)
}
