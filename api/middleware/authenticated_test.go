package middleware_test

//go:generate hel

import (
	"net/http"
	"net/http/httptest"

	"github.com/myshkin5/effective-octo-garbanzo/api/middleware"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authenticated", func() {
	var (
		recorder      *httptest.ResponseRecorder
		request       *http.Request
		mockValidator *mockValidator
		validRequests chan *http.Request
		handler       http.Handler
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		var err error
		request, err = http.NewRequest("GET", "/something", nil)
		Expect(err).NotTo(HaveOccurred())

		mockValidator = newMockValidator()

		validRequests = make(chan *http.Request, 100)

		okFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			validRequests <- r
			w.WriteHeader(http.StatusOK)
		})

		handler = middleware.AuthenticatedHandler(okFunc, "http://auth-server/", mockValidator)
	})

	It("passes the request to the inner handler when the validator says the auth header is valid", func() {
		mockValidator.IsValidOutput.IsValid <- true
		mockValidator.IsValidOutput.Org <- "org1"

		handler.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusOK))
		Expect(mockValidator.IsValidCalled).To(Receive())
		Expect(mockValidator.IsValidInput.AuthHeader).To(Receive(Equal("")))
		var validRequest *http.Request
		Expect(validRequests).To(Receive(&validRequest))
		Expect(validRequest.Context().Value(persistence.OrgContextKey)).To(Equal("org1"))
	})

	It("doesn't pass the request to the inner handler when the validator says the auth header is invalid", func() {
		request.Header.Add("Authorization", "bearer xyz123")
		mockValidator.IsValidOutput.IsValid <- false
		mockValidator.IsValidOutput.Org <- ""

		handler.ServeHTTP(recorder, request)

		Expect(recorder.Code).To(Equal(http.StatusTemporaryRedirect))
		Expect(recorder.Header().Get("Location")).To(Equal("http://auth-server/"))
		Expect(mockValidator.IsValidCalled).To(Receive())
		Expect(mockValidator.IsValidInput.AuthHeader).To(Receive(Equal("bearer xyz123")))
	})
})
