package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cycle-labs/test-app/internal/adapters/api"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthenticatedRouterGroup_WithToken(t *testing.T) {
	testCases := []struct {
		token string
		want  int
	}{
		{"AnyTokenIsOk", http.StatusOK},
		{"", http.StatusUnauthorized},
	}
	for _, tc := range testCases {
		Convey("Given an authenticated route for my gin http server", t, func() {

			recorder := httptest.NewRecorder()
			_, engine := GetTestGinContext(recorder)
			endpoint := api.NewGinAuthenticatorEndpoint()
			group := endpoint.AuthenticatedRouterGroup(engine, "/")
			group.GET("foo", func(*gin.Context) {})

			Convey(fmt.Sprintf("When I perform a get request with authorazation header (Beaer %s)", tc.token), func() {

				req, _ := http.NewRequest("GET", "/foo", nil)
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.token))
				engine.ServeHTTP(recorder, req)

				Convey("The response should have the expected status code", func() {
					got := recorder.Result().StatusCode
					So(got, ShouldEqual, tc.want)
				})
			})
		})
	}
}

func TestAuthenticatedRouterGroup_WithoutToken(t *testing.T) {

	Convey("Given an authenticated route for my gin http server", t, func() {

		recorder := httptest.NewRecorder()
		_, engine := GetTestGinContext(recorder)
		endpoint := api.NewGinAuthenticatorEndpoint()
		group := endpoint.AuthenticatedRouterGroup(engine, "/")
		group.GET("foo", func(*gin.Context) {})

		Convey("When I perform a get request with no authorazation header", func() {
			req, _ := http.NewRequest("GET", "/foo", nil)
			engine.ServeHTTP(recorder, req)

			Convey("The response should have the expected status code", func() {
				So(recorder.Result().StatusCode, ShouldEqual, http.StatusUnauthorized)
			})
		})
	})
}
