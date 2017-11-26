package garbanzo_test

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

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/garbanzo"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

var _ = Describe("Garbanzo", func() {
	const (
		octoName = "kraken"
		url      = "/octos/" + octoName + "/garbanzos/"
	)
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
		garbanzo.MapRoutes("http://here/", router, alice.Chain{}, mockService)
		apiUUID = uuid.NewV4()
	})

	Describe("GET", func() {
		Context("happy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, url+apiUUID.String(), nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchByAPIUUIDAndOctoNameOutput.Garbanzo <- data.Garbanzo{
					APIUUID:      apiUUID,
					GarbanzoType: data.DESI,
					DiameterMM:   4.2,
				}
				mockService.FetchByAPIUUIDAndOctoNameOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("invokes the service layer", func() {
				var actualAPIUUID uuid.UUID
				Expect(mockService.FetchByAPIUUIDAndOctoNameInput.ApiUUID).To(Receive(&actualAPIUUID))
				Expect(actualAPIUUID).To(Equal(apiUUID))
				var actualOctoName string
				Expect(mockService.FetchByAPIUUIDAndOctoNameInput.OctoName).To(Receive(&actualOctoName))
				Expect(actualOctoName).To(Equal(octoName))
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns the garbanzo in the body", func() {
				Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{
					"link":        "http://here%s%s",
					"type":        "DESI",
					"diameter-mm": 4.2
				}`, url, apiUUID)))
			})
		})

		Context("unhappy path", func() {
			Context("invalid UUID", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodGet, url+"not-a-uuid", nil)
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
					request, err = http.NewRequest(http.MethodGet, url+apiUUID.String(), nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.FetchByAPIUUIDAndOctoNameOutput.Garbanzo <- data.Garbanzo{}
					mockService.FetchByAPIUUIDAndOctoNameOutput.Err <- errors.New("bad stuff")

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
					request, err = http.NewRequest(http.MethodGet, url+apiUUID.String(), nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.FetchByAPIUUIDAndOctoNameOutput.Garbanzo <- data.Garbanzo{}
					mockService.FetchByAPIUUIDAndOctoNameOutput.Err <- persistence.ErrNotFound

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
				request, err = http.NewRequest("DELETE", url+apiUUID.String(), nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.DeleteByAPIUUIDAndOctoNameOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("invokes the service layer", func() {
				var actualAPIUUID uuid.UUID
				Expect(mockService.DeleteByAPIUUIDAndOctoNameInput.ApiUUID).To(Receive(&actualAPIUUID))
				Expect(actualAPIUUID).To(Equal(apiUUID))
				var actualOctoName string
				Expect(mockService.DeleteByAPIUUIDAndOctoNameInput.OctoName).To(Receive(&actualOctoName))
				Expect(actualOctoName).To(Equal(octoName))
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
					request, err = http.NewRequest(http.MethodDelete, url+"not-a-uuid", nil)
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
					request, err = http.NewRequest(http.MethodDelete, url+apiUUID.String(), nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.DeleteByAPIUUIDAndOctoNameOutput.Err <- errors.New("bad stuff")

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
					request, err = http.NewRequest(http.MethodDelete, url+apiUUID.String(), nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.DeleteByAPIUUIDAndOctoNameOutput.Err <- persistence.ErrNotFound

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
