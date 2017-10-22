package handlers_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
)

var _ = Describe("Error", func() {
	var (
		recorder *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		recorder.Code = 0

		handlers.Error(recorder, "bad stuff!", http.StatusInternalServerError)
	})

	It("writes the error to a JSON body", func() {
		Expect(recorder.Body).To(MatchJSON(`{
			"code":   500,
			"error":  "bad stuff!",
			"status": "Internal Server Error"
		}`))
	})
})
