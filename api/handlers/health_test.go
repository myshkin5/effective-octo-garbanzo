package handlers_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
)

var _ = Describe("Health", func() {
	var (
		recorder *httptest.ResponseRecorder
		request  *http.Request
		router   *mux.Router
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		router = mux.NewRouter()
		handlers.MapHealthRoutes(router, alice.Chain{})
	})

	Describe("happy path", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest(http.MethodGet, "/health", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a good health body", func() {
			router.ServeHTTP(recorder, request)
			Expect(recorder.Body).To(MatchJSON(`{
				"health": "GOOD"
			}`))
		})

		It("returns an ok status code", func() {
			router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})
	})
})
