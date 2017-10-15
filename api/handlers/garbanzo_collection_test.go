package handlers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

var _ = Describe("GarbanzoCollection", func() {
	var (
		recorder    *httptest.ResponseRecorder
		request     *http.Request
		mockService *mockGarbanzoService
		router      *mux.Router
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		mockService = newMockGarbanzoService()

		router = mux.NewRouter()
		handlers.MapGarbanzoCollectionRoutes("http://here/", router, alice.Chain{}, mockService)
	})

	Describe("GET", func() {
		Context("happy path - empty collection", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/garbanzos", nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchAllGarbanzosOutput.Garbanzos <- []persistence.Garbanzo{}
				mockService.FetchAllGarbanzosOutput.Err <- nil
			})

			It("returns an ok status code", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns an empty list in the body", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Body).To(MatchJSON(`{
					"data": {
						"garbanzos": []
					}
				}`))
			})
		})

		Context("happy path", func() {
			var (
				apiUUID1 uuid.UUID
				apiUUID2 uuid.UUID
			)

			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/garbanzos", nil)
				Expect(err).NotTo(HaveOccurred())

				apiUUID1 = uuid.NewV4()
				apiUUID2 = uuid.NewV4()

				mockService.FetchAllGarbanzosOutput.Garbanzos <- []persistence.Garbanzo{
					{
						APIUUID:   apiUUID1,
						FirstName: "Joe",
						LastName:  "Schmoe",
					},
					{
						APIUUID:   apiUUID2,
						FirstName: "Marty",
						LastName:  "Blarty",
					},
				}
				mockService.FetchAllGarbanzosOutput.Err <- nil
			})

			It("returns an ok status code", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns all garbanzos in the body", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{
					"data": {
						"garbanzos": [
							{
								"link": "http://here/garbanzos/%s",
								"first-name": "Joe",
								"last-name": "Schmoe"
							},
							{
								"link": "http://here/garbanzos/%s",
								"first-name": "Marty",
								"last-name": "Blarty"
							}
						]
					}
				}`, apiUUID1, apiUUID2)))
			})
		})

		Context("unhappy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/garbanzos", nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchAllGarbanzosOutput.Garbanzos <- nil
				mockService.FetchAllGarbanzosOutput.Err <- errors.New("bad stuff")
			})

			It("returns an internal server error status code", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})

			It("returns a JSON error", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Body).To(MatchJSON(`{
					"code": 500,
					"error": "Error fetching all garbanzos",
					"status": "Internal Server Error"
				}`))
			})
		})
	})

	Describe("POST", func() {
		Context("happy path", func() {
			var (
				apiUUID uuid.UUID
			)

			BeforeEach(func() {
				var err error
				body := strings.NewReader(`{
					"first-name": "Joe",
					"last-name":  "Schmoe",
					"something":  "ignored"
				}`)
				request, err = http.NewRequest(http.MethodPost, "/garbanzos", body)
				Expect(err).NotTo(HaveOccurred())

				apiUUID = uuid.NewV4()

				mockService.CreateGarbanzoOutput.GarbanzoOut <- persistence.Garbanzo{
					Id:        234,
					APIUUID:   apiUUID,
					FirstName: "Joe",
					LastName:  "Schmoe",
				}
				mockService.CreateGarbanzoOutput.Err <- nil
			})

			It("returns an ok status code", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})

			It("returns the newly created garbanzo in the body", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{
					"data": {
						"garbanzo": {
							"link": "http://here/garbanzos/%s",
							"first-name": "Joe",
							"last-name": "Schmoe"
						}
					}
				}`, apiUUID)))
			})
		})

		Context("unhappy path", func() {
			Context("invalid json", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader("not json")
					request, err = http.NewRequest(http.MethodPost, "/garbanzos", body)
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns an internal server error status code", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 400,
						"error": "Body of request was not valid JSON",
						"status": "Bad Request"
					}`))
				})
			})

			Context("persistence error", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader("{}")
					request, err = http.NewRequest(http.MethodPost, "/garbanzos", body)
					Expect(err).NotTo(HaveOccurred())

					mockService.CreateGarbanzoOutput.GarbanzoOut <- persistence.Garbanzo{}
					mockService.CreateGarbanzoOutput.Err <- errors.New("not good")
				})

				It("returns an internal server error status code", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 500,
						"error": "Error creating new garbanzo",
						"status": "Internal Server Error"
					}`))
				})
			})
		})
	})
})
