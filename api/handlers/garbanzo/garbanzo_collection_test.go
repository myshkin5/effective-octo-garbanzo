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
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
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
		garbanzo.MapCollectionRoutes("http://here/", router, alice.Chain{}, mockService)
	})

	Describe("GET", func() {
		Context("happy path - empty collection", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/garbanzos", nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchAllGarbanzosOutput.Garbanzos <- []data.Garbanzo{}
				mockService.FetchAllGarbanzosOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns an empty list in the body", func() {
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

				mockService.FetchAllGarbanzosOutput.Garbanzos <- []data.Garbanzo{
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
				mockService.FetchAllGarbanzosOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns all garbanzos in the body", func() {
				Expect(recorder.Body).To(MatchJSON(fmt.Sprintf(`{
					"data": {
						"garbanzos": [
							{
								"link":        "http://here/garbanzos/%s",
								"type":        "DESI",
								"diameter-mm": 4.2
							},
							{
								"link":        "http://here/garbanzos/%s",
								"type":        "KABULI",
								"diameter-mm": 6.4
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
				request, err = http.NewRequest(http.MethodPost, "/garbanzos", body)
				Expect(err).NotTo(HaveOccurred())

				apiUUID = uuid.NewV4()

				mockService.CreateGarbanzoOutput.GarbanzoOut <- data.Garbanzo{
					Id:           234,
					APIUUID:      apiUUID,
					GarbanzoType: data.DESI,
					DiameterMM:   4.2,
				}
				mockService.CreateGarbanzoOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("creates the garbanzo via the service", func() {
				var garbanzo data.Garbanzo
				Expect(mockService.CreateGarbanzoInput.GarbanzoIn).To(Receive(&garbanzo))
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
			Context("invalid json", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader("not json")
					request, err = http.NewRequest(http.MethodPost, "/garbanzos", body)
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
					request, err = http.NewRequest(http.MethodPost, "/garbanzos", body)
					Expect(err).NotTo(HaveOccurred())

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 400,
						"error": "Unknown garbanzo type: RED",
						"status": "Bad Request"
					}`))
				})
			})

			Context("persistence error", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader(`{
						"type":        "DESI",
						"diameter-mm": 4.2
					}`)
					request, err = http.NewRequest(http.MethodPost, "/garbanzos", body)
					Expect(err).NotTo(HaveOccurred())

					mockService.CreateGarbanzoOutput.GarbanzoOut <- data.Garbanzo{}
					mockService.CreateGarbanzoOutput.Err <- errors.New("not good")

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
		})
	})
})
