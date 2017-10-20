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
		handlers.MapCatchAllRoutes("http://here/", router, alice.Chain{})
	})

	Context("happy path", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest(http.MethodGet, "/", nil)
			request.RequestURI = "/"
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns a good root body", func() {
			router.ServeHTTP(recorder, request)
			Expect(recorder.Body).To(MatchJSON(`{
				"health": "http://here/health",
				"garbanzos": "http://here/garbanzos"
			}`))
		})

		It("returns an ok status code", func() {
			router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})
	})

	Context("catch all", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest(http.MethodGet, "/any-thing", nil)
			request.RequestURI = "/any-thing"
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
