package middleware_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/middleware"
)

var _ = Describe("StandardHeaders", func() {
	var (
		recorder *httptest.ResponseRecorder
		request  *http.Request
		handler  http.Handler
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		var err error
		request, err = http.NewRequest("GET", "/something", nil)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("200, Ok", func() {
		BeforeEach(func() {
			handler = middleware.StandardHeadersHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
		})

		It("adds the standard header", func() {
			handler.ServeHTTP(recorder, request)

			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusOK))
		})
	})

	Context("no content status code", func() {
		BeforeEach(func() {
			handler = middleware.StandardHeadersHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
		})

		It("does not add the standard header", func() {
			handler.ServeHTTP(recorder, request)

			Expect(recorder.Header().Get("Content-Type")).To(Equal(""))

			Expect(recorder.Code).To(Equal(http.StatusNoContent))
		})
	})

	Describe("writes data when the header has been written", func() {
		BeforeEach(func() {
			handler = middleware.StandardHeadersHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{ "heres": "some-json" }`))
			}))
		})

		It("adds the standard header", func() {
			handler.ServeHTTP(recorder, request)

			Expect(recorder.Body).To(MatchJSON(`{ "heres": "some-json" }`))
		})
	})
})
