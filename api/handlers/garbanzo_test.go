package handlers_test

//go:generate hel

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

var _ = Describe("Garbanzo", func() {
	var (
		recorder    *httptest.ResponseRecorder
		request     *http.Request
		mockService *mockGarbanzoService
		router      *mux.Router
		apiUUID     uuid.UUID
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		mockService = newMockGarbanzoService()

		router = mux.NewRouter()
		handlers.MapGarbanzoRoutes("http://here/", router, alice.Chain{}, mockService)
		apiUUID = uuid.NewV4()
	})

	Describe("GET", func() {
		Context("happy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/garbanzos/"+apiUUID.String(), nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchGarbanzoByAPIUUIDOutput.Garbanzo <- persistence.Garbanzo{
					APIUUID:      apiUUID,
					GarbanzoType: persistence.DESI,
					DiameterMM:   4.2,
				}
				mockService.FetchGarbanzoByAPIUUIDOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns the garbanzo in the body", func() {
				Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{
					"data": {
						"garbanzo": {
							"link":        "http://here/garbanzos/%s",
							"type":        "DESI",
							"diameter-mm": 4.2
						}
					}
				}`, apiUUID)))
			})
		})

		Context("unhappy path", func() {
			Context("invalid UUID", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodGet, "/garbanzos/not-a-uuid", nil)
					Expect(err).NotTo(HaveOccurred())

					router.ServeHTTP(recorder, request)
				})

				It("returns a bad request status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 400,
						"error": "Invalid UUID",
						"status": "Bad Request"
					}`))
				})
			})

			Context("persistence error", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodGet, "/garbanzos/"+apiUUID.String(), nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.FetchGarbanzoByAPIUUIDOutput.Garbanzo <- persistence.Garbanzo{}
					mockService.FetchGarbanzoByAPIUUIDOutput.Err <- errors.New("bad stuff")

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 500,
						"error": "Error fetching garbanzo",
						"status": "Internal Server Error"
					}`))
				})
			})

			Context("not found error", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodGet, "/garbanzos/"+apiUUID.String(), nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.FetchGarbanzoByAPIUUIDOutput.Garbanzo <- persistence.Garbanzo{}
					mockService.FetchGarbanzoByAPIUUIDOutput.Err <- persistence.ErrNotFound

					router.ServeHTTP(recorder, request)
				})

				It("returns a not found status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{
						"code": 404,
						"error": "Garbanzo %s not found",
						"status": "Not Found"
					}`, apiUUID)))
				})
			})
		})
	})

	Describe("DELETE", func() {
		Context("happy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest("DELETE", "/garbanzos/"+apiUUID.String(), nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.DeleteGarbanzoByAPIUUIDOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns a no content status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusNoContent))
			})

			It("returns no content", func() {
				Expect(recorder.Body.Len()).To(BeZero())
			})
		})

		Context("unhappy path", func() {
			Context("invalid UUID", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodDelete, "/garbanzos/not-a-uuid", nil)
					Expect(err).NotTo(HaveOccurred())

					router.ServeHTTP(recorder, request)
				})

				It("returns a bad request status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 400,
						"error": "Invalid UUID",
						"status": "Bad Request"
					}`))
				})
			})

			Context("persistence error", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodDelete, "/garbanzos/"+apiUUID.String(), nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.DeleteGarbanzoByAPIUUIDOutput.Err <- errors.New("bad stuff")

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 500,
						"error": "Error fetching garbanzo",
						"status": "Internal Server Error"
					}`))
				})
			})

			Context("not found error", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodDelete, "/garbanzos/"+apiUUID.String(), nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.DeleteGarbanzoByAPIUUIDOutput.Err <- persistence.ErrNotFound

					router.ServeHTTP(recorder, request)
				})

				It("returns a not found status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{
						"code": 404,
						"error": "Garbanzo %s not found",
						"status": "Not Found"
					}`, apiUUID)))
				})
			})
		})
	})
})
