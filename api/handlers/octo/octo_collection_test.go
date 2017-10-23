package octo_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/octo"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

var _ = Describe("OctoCollection", func() {
	var (
		recorder    *httptest.ResponseRecorder
		request     *http.Request
		mockService *mockOctoService
		router      *mux.Router
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		mockService = newMockOctoService()

		router = mux.NewRouter()
		octo.MapCollectionRoutes("http://here/", router, alice.Chain{}, mockService)
	})

	Describe("GET", func() {
		Context("happy path - empty collection", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/octos", nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchAllOctosOutput.Octos <- []data.Octo{}
				mockService.FetchAllOctosOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns an empty list in the body", func() {
				Expect(recorder.Body).To(MatchJSON(`{
					"data": {
						"octos": []
					}
				}`))
			})
		})

		Context("happy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/octos", nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchAllOctosOutput.Octos <- []data.Octo{
					{
						Name: "kraken",
					},
					{
						Name: "cthulhu",
					},
				}
				mockService.FetchAllOctosOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns all octos in the body", func() {
				Expect(recorder.Body).To(MatchJSON(`{
					"data": {
						"octos": [
							{
								"link":        "http://here/octos/kraken",
								"name":        "kraken"
							},
							{
								"link":        "http://here/octos/cthulhu",
								"name":        "cthulhu"
							}
						]
					}
				}`))
			})
		})

		Context("unhappy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/octos", nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchAllOctosOutput.Octos <- nil
				mockService.FetchAllOctosOutput.Err <- errors.New("bad stuff")

				router.ServeHTTP(recorder, request)
			})

			It("returns an internal server error status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			})

			It("returns a JSON error", func() {
				Expect(recorder.Body).To(MatchJSON(`{
					"code": 500,
					"error": "Error fetching all octos",
					"status": "Internal Server Error"
				}`))
			})
		})
	})

	Describe("POST", func() {
		Context("happy path", func() {
			BeforeEach(func() {
				var err error
				body := strings.NewReader(`{
					"name": "kraken"
				}`)
				request, err = http.NewRequest(http.MethodPost, "/octos", body)
				Expect(err).NotTo(HaveOccurred())

				mockService.CreateOctoOutput.OctoOut <- data.Octo{
					Id:   234,
					Name: "kraken",
				}
				mockService.CreateOctoOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("creates the octo via the service", func() {
				var octo data.Octo
				Expect(mockService.CreateOctoInput.OctoIn).To(Receive(&octo))
				Expect(octo).To(Equal(data.Octo{
					Name: "kraken",
				}))
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})

			It("returns the newly created octo in the body", func() {
				Expect(recorder.Body).To(MatchJSON(`{
					"data": {
						"octo": {
							"link": "http://here/octos/kraken",
							"name": "kraken"
						}
					}
				}`))
			})
		})

		Context("unhappy path", func() {
			Context("invalid json", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader("not json")
					request, err = http.NewRequest(http.MethodPost, "/octos", body)
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

			Context("persistence error", func() {
				BeforeEach(func() {
					var err error
					body := strings.NewReader(`{
						"name": "kraken"
					}`)
					request, err = http.NewRequest(http.MethodPost, "/octos", body)
					Expect(err).NotTo(HaveOccurred())

					mockService.CreateOctoOutput.OctoOut <- data.Octo{}
					mockService.CreateOctoOutput.Err <- errors.New("not good")

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 500,
						"error": "Error creating new octo",
						"status": "Internal Server Error"
					}`))
				})
			})
		})
	})
})
