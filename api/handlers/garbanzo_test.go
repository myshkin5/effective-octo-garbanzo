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
					APIUUID:   apiUUID,
					FirstName: "Joe",
					LastName:  "Schmoe",
				}
				mockService.FetchGarbanzoByAPIUUIDOutput.Err <- nil
			})

			It("returns an ok status code", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns the garbanzo in the body", func() {
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
			Context("invalid UUID", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodGet, "/garbanzos/not-a-uuid", nil)
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns a bad request status code", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error", func() {
					router.ServeHTTP(recorder, request)
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
				})

				It("returns an internal server error status code", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					router.ServeHTTP(recorder, request)
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
				})

				It("returns a not found status code", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				})

				It("returns a JSON error", func() {
					router.ServeHTTP(recorder, request)
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
			})

			It("returns a no content status code", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusNoContent))
			})

			It("returns no content", func() {
				router.ServeHTTP(recorder, request)
				Expect(recorder.Body.Len()).To(BeZero())
			})
		})

		Context("unhappy path", func() {
			Context("invalid UUID", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodDelete, "/garbanzos/not-a-uuid", nil)
					Expect(err).NotTo(HaveOccurred())
				})

				It("returns a bad request status code", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error", func() {
					router.ServeHTTP(recorder, request)
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
				})

				It("returns an internal server error status code", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					router.ServeHTTP(recorder, request)
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
				})

				It("returns a not found status code", func() {
					router.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				})

				It("returns a JSON error", func() {
					router.ServeHTTP(recorder, request)
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
