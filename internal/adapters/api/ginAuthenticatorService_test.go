package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthenticatedToken(t *testing.T) {
	testCases := []struct {
		token string
		want  int
	}{
		{"AnyTokenIsOk", http.StatusOK},
		{"", http.StatusUnauthorized},
	}
	for _, tc := range testCases {
		recorder := httptest.NewRecorder()
		ctx := GetTestGinContext(recorder)
		MockJsonGet(ctx, WithBearerToken(tc.token))

		service := NewGinAuthenicatorService()
		service.AuthenticateToken(ctx)

		got := recorder.Result().StatusCode
		if got != tc.want {
			t.Errorf("Want: %d Got: %d", tc.want, got)
		}
	}
}

func TestUnauthenticatedToken(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := GetTestGinContext(recorder)
	MockJsonGet(ctx)

	service := NewGinAuthenicatorService()
	service.AuthenticateToken(ctx)

	want := http.StatusUnauthorized
	got := recorder.Result().StatusCode
	if got != want {
		t.Errorf("Want: %d Got: %d", want, got)
	}
}

type Option = func(ctx *gin.Context)

func WithBearerToken(token string) Option {
	return func(c *gin.Context) {
		c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

func WithPathParams(pathParams gin.Params) Option {
	return func(c *gin.Context) {
		c.Params = pathParams
	}
}

func WithQueryParams(queryParams url.Values) Option {
	return func(c *gin.Context) {
		c.Request.URL.RawQuery = queryParams.Encode()
	}
}

func MockJsonGet(c *gin.Context, opts ...Option) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")
	for _, opt := range opts {
		opt(c)
	}
}

func GetTestGinContext(recorder *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	return ctx
}
