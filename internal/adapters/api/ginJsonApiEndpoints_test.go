package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dumpsterfireproject/test-app/internal/adapters/api"
	"github.com/dumpsterfireproject/test-app/internal/adapters/graph/model"
	"github.com/genjidb/genji"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

var numberOfSeedItems = 5

type TestJsonServer struct {
	Recorder *httptest.ResponseRecorder
	Engine   *gin.Engine
	Context  *gin.Context
	DB       *genji.DB
}

func NewTestJsonServer(t *testing.T) TestJsonServer {
	recorder := httptest.NewRecorder()
	ctx, engine := GetTestGinContext(recorder)
	endpoint := api.NewGinAuthenticatorEndpoint()
	group := endpoint.AuthenticatedRouterGroup(engine, "/")
	db, service := InMemoryStockService(t)
	endpoints := api.NewGinJsonAPIEndpoints(service)
	endpoints.AddInventoryRoutes(group)
	return TestJsonServer{
		Recorder: recorder,
		Engine:   engine,
		Context:  ctx,
		DB:       db,
	}
}

func (s TestJsonServer) ServeHTTP(req *http.Request) {
	s.Engine.ServeHTTP(s.Recorder, req)
}

func (s TestJsonServer) Close() {
	s.DB.Close()
}

func authenticatedRequest(method string, path string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, path, body)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))
	return req
}

func TestAddInventoryRoutes_GetItems(t *testing.T) {
	testCases := []struct {
		wantStatus int
		wantBody   string
	}{
		{http.StatusOK, ""},
	}
	for _, tc := range testCases {

		Convey("Given a JSON HTTP adapter to access an InventoryStockService", t, func() {
			testServer := NewTestJsonServer(t)
			defer testServer.Close()

			Convey("When an authorized GET request is made to /items", func() {
				req := authenticatedRequest("GET", "/items", nil)
				testServer.ServeHTTP(req)

				Convey("Then response status code should be 200", func() {
					got := testServer.Recorder.Result().StatusCode
					So(got, ShouldEqual, tc.wantStatus)
				})

				Convey("And the expected number of items should be returned", func() {
					var got []*model.Item
					err := json.Unmarshal([]byte(testServer.Recorder.Body.String()), &got)
					So(err, ShouldBeNil)
					So(len(got), ShouldEqual, numberOfSeedItems)
				})
			})

		})
	}
}

func TestAddInventoryRoutes_GetInventory(t *testing.T) {
	testCases := []struct {
		path               string
		wantStatus         int
		wantElementsInBody int
	}{
		// wantElementsInBody based on seed data
		{"/inventory", http.StatusOK, 6},
		{"/inventory?itemID=8e1af20d-7c39-47e2-a70c-3938bcee2e29", http.StatusOK, 2},
	}
	for _, tc := range testCases {

		Convey("Given a JSON HTTP adapter to access an InventoryStockService", t, func() {
			testServer := NewTestJsonServer(t)
			defer testServer.Close()

			Convey("When an authorized GET request is made", func() {
				req := authenticatedRequest("GET", tc.path, nil)
				testServer.ServeHTTP(req)

				Convey("Then response status code should be 200", func() {
					got := testServer.Recorder.Result().StatusCode
					SoMsg(tc.path, got, ShouldEqual, tc.wantStatus)
				})

				Convey("And the expected number of inventory records should be returned", func() {
					var got []*model.Inventory
					err := json.Unmarshal([]byte(testServer.Recorder.Body.String()), &got)
					SoMsg(tc.path, err, ShouldBeNil)
					SoMsg(tc.path, len(got), ShouldEqual, tc.wantElementsInBody)
				})
			})

		})
	}
}

func TestAddInventoryRoutes_AddInventory(t *testing.T) {
	testCases := []struct {
		path       string
		wantStatus int
	}{
		// wantElementsInBody based on seed data
		{"/addInventory", http.StatusCreated},
	}
	for _, tc := range testCases {

		Convey("Given a JSON HTTP adapter to access an InventoryStockService", t, func() {
			testServer := NewTestJsonServer(t)
			defer testServer.Close()

			Convey("When an authorized POST request is made", func() {
				inventory := &model.Inventory{
					ID: "",
					Item: &model.Item{
						ID:          "8e1af20d-7c39-47e2-a70c-3938bcee2e29",
						Sku:         "60760-3400",
						Description: "Modafinil",
					},
					Location: "Memphis",
					Status:   "Quarantine",
					Quantity: 100,
				}
				body, _ := json.Marshal(inventory)
				req := authenticatedRequest("POST", tc.path, bytes.NewBuffer(body))
				testServer.ServeHTTP(req)

				Convey("Then response status code should be 201", func() {
					got := testServer.Recorder.Result().StatusCode
					SoMsg(tc.path, got, ShouldEqual, tc.wantStatus)
				})
			})

		})
	}
}
