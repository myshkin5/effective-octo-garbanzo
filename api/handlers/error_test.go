package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/services"
)

var _ = Describe("Error", func() {
	var (
		recorder *httptest.ResponseRecorder
	)

	Context("regular errors", func() {
		BeforeEach(func() {
			recorder = httptest.NewRecorder()
			recorder.Code = 0

			handlers.Error(recorder, "bad stuff!", http.StatusInternalServerError, nil, nil)
		})

		It("writes the error to a JSON body", func() {
			Expect(recorder.Body).To(MatchJSON(`{
				"code":   500,
				"error":  "bad stuff!",
				"status": "Internal Server Error"
			}`))
		})
	})

	Context("validation errors", func() {
		BeforeEach(func() {
			recorder = httptest.NewRecorder()
			recorder.Code = 0

			err := services.NewValidationError(map[string][]string{
				"FieldA": {"1", "2"},
				"FieldB": {"3", "4"},
				"FieldC": {"5", "6"},
			})
			mapping := map[string]string{"FieldA": "field-a", "FieldB": "field-b"}
			handlers.Error(recorder, "bad stuff!", http.StatusInternalServerError, err, mapping)
		})

		It("writes the error to a JSON body", func() {
			/* The errors are in an unordered array so we can't just MatchJSON(). The full doc should look like this:
			{
				"code": 400,
				"error": "bad stuff!",
				"errors": [
					"field-a 1",
					"field-a 2",
					"field-b 3",
					"field-b 4",
					"FieldC 5",
					"FieldC 6"
				],
				"status": "Bad Request"
			}
			*/
			jsonObj := handlers.JSONObject{}
			err := json.NewDecoder(recorder.Body).Decode(&jsonObj)
			Expect(err).NotTo(HaveOccurred())

			Expect(jsonObj["code"]).To(BeEquivalentTo(400))
			Expect(jsonObj["error"]).To(Equal("bad stuff!"))
			Expect(jsonObj["status"]).To(Equal("Bad Request"))

			errors, ok := jsonObj["errors"].([]interface{})
			Expect(ok).To(BeTrue())
			Expect(errors).To(ContainElement("field-a 1"))
			Expect(errors).To(ContainElement("field-a 2"))
			Expect(errors).To(ContainElement("field-b 3"))
			Expect(errors).To(ContainElement("field-b 4"))
			Expect(errors).To(ContainElement("FieldC 5"))
			Expect(errors).To(ContainElement("FieldC 6"))
		})
	})
})
