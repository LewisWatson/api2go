// +build gingonic,!gorillamux,!echo

package routing_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/LewisWatson/api2go"
	"github.com/LewisWatson/api2go/examples/model"
	"github.com/LewisWatson/api2go/examples/resource"
	"github.com/LewisWatson/api2go/examples/storage"
	"github.com/LewisWatson/api2go/routing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("api2go with gingonic router adapter", func() {
	var (
		router routing.Routeable
		gg     *gin.Engine
		api    *api2go.API
		rec    *httptest.ResponseRecorder
	)

	BeforeSuite(func() {
		gin.SetMode(gin.ReleaseMode)
		gg = gin.Default()
		router = routing.Gin(gg)
		api = api2go.NewAPIWithRouting(
			"api",
			api2go.NewStaticResolver("/"),
			router,
		)

		userStorage := storage.NewUserStorage()
		chocStorage := storage.NewChocolateStorage()
		api.AddResource(model.User{}, resource.UserResource{ChocStorage: chocStorage, UserStorage: userStorage})
		api.AddResource(model.Chocolate{}, resource.ChocolateResource{ChocStorage: chocStorage, UserStorage: userStorage})
	})

	BeforeEach(func() {
		log.SetOutput(ioutil.Discard)
		rec = httptest.NewRecorder()
	})

	Context("CRUD Tests", func() {
		It("will create a new user", func() {
			reqBody := strings.NewReader(`{"data": {"attributes": {"user-name": "Sansa Stark"}, "id": "1", "type": "users"}}`)
			req, err := http.NewRequest("POST", "/api/users", reqBody)
			Expect(err).To(BeNil())
			gg.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusCreated))
		})

		It("will find her", func() {
			expectedUser := `
			{
				"data":
				{
					"attributes":{
						"user-name":"Sansa Stark"
					},
					"id":"1",
					"relationships":{
						"sweets":{
							"data":[],"links":{"related":"/api/users/1/sweets","self":"/api/users/1/relationships/sweets"}
						}
					},"type":"users"
				},
				"meta":
				{
					"author":"The api2go examples crew","license":"wtfpl","license-url":"http://www.wtfpl.net"
				}
			}`

			req, err := http.NewRequest("GET", "/api/users/1", nil)
			Expect(err).To(BeNil())
			gg.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(string(rec.Body.Bytes())).To(MatchJSON((expectedUser)))
		})

		It("can call handle", func() {
			handler := api.Handler()
			_, ok := handler.(http.Handler)
			Expect(ok).To(Equal(true))
		})

		It("update the username", func() {
			reqBody := strings.NewReader(`{"data": {"id": "1", "attributes": {"user-name": "Alayne"}, "type" : "users"}}`)
			req, err := http.NewRequest("PATCH", "/api/users/1", reqBody)
			Expect(err).To(BeNil())
			gg.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNoContent))
		})

		It("will find her once again", func() {
			expectedUser := `
			{
				"data":
				{
					"attributes":{
						"user-name":"Alayne"
					},
					"id":"1",
					"relationships":{
						"sweets":{
							"data":[],"links":{"related":"/api/users/1/sweets","self":"/api/users/1/relationships/sweets"}
						}
					},"type":"users"
				},
				"meta":
				{
					"author":"The api2go examples crew","license":"wtfpl","license-url":"http://www.wtfpl.net"
				}
			}`

			req, err := http.NewRequest("GET", "/api/users/1", nil)
			Expect(err).To(BeNil())
			gg.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(string(rec.Body.Bytes())).To(MatchJSON((expectedUser)))
		})

		It("will delete her", func() {
			req, err := http.NewRequest("DELETE", "/api/users/1", nil)
			Expect(err).To(BeNil())
			gg.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNoContent))
		})

		It("won't find her anymore", func() {
			expected := `{"errors":[{"status":"404","title":"http error (404) User for id 1 not found and 0 more errors, User for id 1 not found"}]}`
			req, err := http.NewRequest("GET", "/api/users/1", nil)
			Expect(err).To(BeNil())
			gg.ServeHTTP(rec, req)
			Expect(rec.Code).To(Equal(http.StatusNotFound))
			Expect(string(rec.Body.Bytes())).To(MatchJSON(expected))
		})
	})
})
