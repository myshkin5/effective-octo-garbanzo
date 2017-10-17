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

var _ = Describe("CatchAll", func() {
	var (
		recorder *httptest.ResponseRecorder
		request  *http.Request
		router   *mux.Router
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		router = mux.NewRouter()
		handlers.MapCatchAllRoutes(router, alice.Chain{})
	})

	Describe("catch all", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest(http.MethodGet, "/", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error body", func() {
			router.ServeHTTP(recorder, request)
			Expect(recorder.Body).To(MatchJSON(`{
				"code": 404,
				"error": "Not Found",
				"status": "Not Found"
			}`))
		})

		It("returns a not found code", func() {
			router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})
	})
})
