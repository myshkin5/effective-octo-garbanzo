package octo_test

//go:generate hel

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/octo"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

var _ = Describe("Octo", func() {
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
		octo.MapRoutes("http://here/", router, alice.Chain{}, mockService)
	})

	Describe("GET", func() {
		Context("happy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest(http.MethodGet, "/octos/kraken", nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.FetchOctoByNameOutput.Octo <- data.Octo{
					Name: "kraken",
				}
				mockService.FetchOctoByNameOutput.Err <- nil

				router.ServeHTTP(recorder, request)
			})

			It("returns an ok status code", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("returns the octo in the body", func() {
				Expect(recorder.Body).To(MatchJSON(`{
					"link":      "http://here/octos/kraken",
					"name":      "kraken",
					"garbanzos": "http://here/octos/kraken/garbanzos"
				}`))
			})
		})

		Context("unhappy path", func() {
			Context("persistence error", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodGet, "/octos/kraken", nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.FetchOctoByNameOutput.Octo <- data.Octo{}
					mockService.FetchOctoByNameOutput.Err <- errors.New("bad stuff")

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 500,
						"error": "Error fetching octo",
						"status": "Internal Server Error"
					}`))
				})
			})

			Context("not found error", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodGet, "/octos/squidward", nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.FetchOctoByNameOutput.Octo <- data.Octo{}
					mockService.FetchOctoByNameOutput.Err <- persistence.ErrNotFound

					router.ServeHTTP(recorder, request)
				})

				It("returns a not found status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 404,
						"error": "Octo squidward not found",
						"status": "Not Found"
					}`))
				})
			})
		})
	})

	Describe("DELETE", func() {
		Context("happy path", func() {
			BeforeEach(func() {
				var err error
				request, err = http.NewRequest("DELETE", "/octos/kraken", nil)
				Expect(err).NotTo(HaveOccurred())

				mockService.DeleteOctoByNameOutput.Err <- nil

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
			Context("persistence error", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodDelete, "/octos/kraken", nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.DeleteOctoByNameOutput.Err <- errors.New("bad stuff")

					router.ServeHTTP(recorder, request)
				})

				It("returns an internal server error status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 500,
						"error": "Error fetching octo",
						"status": "Internal Server Error"
					}`))
				})
			})

			Context("not found error", func() {
				BeforeEach(func() {
					var err error
					request, err = http.NewRequest(http.MethodDelete, "/octos/squidward", nil)
					Expect(err).NotTo(HaveOccurred())

					mockService.DeleteOctoByNameOutput.Err <- persistence.ErrNotFound

					router.ServeHTTP(recorder, request)
				})

				It("returns a not found status code", func() {
					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				})

				It("returns a JSON error", func() {
					Expect(recorder.Body).To(MatchJSON(`{
						"code": 404,
						"error": "Octo squidward not found",
						"status": "Not Found"
					}`))
				})
			})
		})
	})
})
