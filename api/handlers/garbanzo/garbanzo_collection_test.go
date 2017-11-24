package garbanzo_test

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

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/garbanzo"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

var _ = Describe("GarbanzoCollection", func() {
	const url = "/octos/kraken/garbanzos"
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
		garbanzo.MapCollectionRoutes("http://here/", router, alice.Chain{}, mockService)
	})

	Describe("GET", func() {
		Context("happy path - empty collection", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, url, nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchAllOutput.Garbanzos <- []data.Garbanzo{}
				mockService.FetchAllOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns an empty list in the body", func() {
				Expect(recorder.Body).To(MatchJSON(`[]`))
			})
		})

		Context("happy path", func() {
			var (
				apiUUID1 uuid.UUID
				apiUUID2 uuid.UUID
			)

			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, url, nil)
				Expect(err).NotTo(HaveOccurred())

				apiUUID1 = uuid.NewV4()
				apiUUID2 = uuid.NewV4()

				mockService.FetchAllOutput.Garbanzos <- []data.Garbanzo{
					{
						APIUUID:      apiUUID1,
						GarbanzoType: data.DESI,
						DiameterMM:   4.2,
					},
					{
						APIUUID:      apiUUID2,
						GarbanzoType: data.KABULI,
						DiameterMM:   6.4,
					},
				}
				mockService.FetchAllOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns all garbanzos in the body", func() {
				Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`[
					{
						"link":        "http://here%s/%s",
						"type":        "DESI",
						"diameter-mm": 4.2
					},
					{
						"link":        "http://here%s/%s",
						"type":        "KABULI",
						"diameter-mm": 6.4
					}
				]`, url, apiUUID1, url, apiUUID2)))
			})
		})

		Context("unhappy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, url, nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchAllOutput.Garbanzos <- nil
				mockService.FetchAllOutput.Err <- errors.New("bad stuff")

				router.ServeHTTP(recorder, request)
			})

			It("returns an internal server error status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})

			It("returns a JSON error", func() {
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
					"type":        "DESI",
					"diameter-mm": 4.2,
					"something":   "ignored"
				}`)
				request, err = http.NewRequest(http.MethodPost, url, body)
				Expect(err).NotTo(HaveOccurred())

				apiUUID = uuid.NewV4()

				mockService.CreateOutput.GarbanzoOut <- data.Garbanzo{
					Id:           234,
					APIUUID:      apiUUID,
					GarbanzoType: data.DESI,
					DiameterMM:   4.2,
				}
				mockService.CreateOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("creates the garbanzo via the service", func() {
				Expect(mockService.CreateCalled).To(HaveLen(1))
				var octoName string
				Expect(mockService.CreateInput.OctoName).To(Receive(&octoName))
				Expect(octoName).To(Equal("kraken"))
				var garbanzo data.Garbanzo
				Expect(mockService.CreateInput.GarbanzoIn).To(Receive(&garbanzo))
				Expect(garbanzo).To(Equal(data.Garbanzo{
					GarbanzoType: data.DESI,
					DiameterMM:   4.2,
				}))
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})

			It("returns the newly created garbanzo in the body", func() {
				Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{
					"link":        "http://here%s/%s",
					"type":        "DESI",
					"diameter-mm": 4.2
				}`, url, apiUUID)))
			})
		})

		Context("unhappy path", func() {
			Context("invalid json", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader("not json")
					request, err = http.NewRequest(http.MethodPost, url, body)
					Expect(err).NotTo(HaveOccurred())

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 400,
						"error": "Body of request was not valid JSON",
						"status": "Bad Request"
					}`))
				})
			})

			Context("invalid type", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader(`{
						"type":        "RED",
						"diameter-mm": 4.2
					}`)
					request, err = http.NewRequest(http.MethodPost, url, body)
					Expect(err).NotTo(HaveOccurred())

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 400,
						"error": "invalid garbanzo type: RED",
						"status": "Bad Request"
					}`))
				})
			})

			Context("general persistence error", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader(`{
						"type":        "DESI",
						"diameter-mm": 4.2
					}`)
					request, err = http.NewRequest(http.MethodPost, url, body)
					Expect(err).NotTo(HaveOccurred())

					mockService.CreateOutput.GarbanzoOut <- data.Garbanzo{}
					mockService.CreateOutput.Err <- errors.New("not good")

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 500,
						"error": "Error creating new garbanzo",
						"status": "Internal Server Error"
					}`))
				})
			})

			Context("not found persistence error", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader(`{
						"type":        "DESI",
						"diameter-mm": 4.2
					}`)
					request, err = http.NewRequest(http.MethodPost, url, body)
					Expect(err).NotTo(HaveOccurred())

					mockService.CreateOutput.GarbanzoOut <- data.Garbanzo{}
					mockService.CreateOutput.Err <- persistence.ErrNotFound

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusConflict))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 409,
						"error": "Parent octo 'kraken' not found",
						"status": "Conflict"
					}`))
				})
			})
		})
	})
})
