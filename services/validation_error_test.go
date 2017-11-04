package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/services"
)

var _ = Describe("ValidationError", func() {
	It("returns a nice error string", func() {
		errors := make(map[string][]string)
		errors["field1"] = []string{"must be present", "must be blue or chartreuse"}
		errors["field2"] = []string{"must be unhinged"}
		err := services.NewValidationError(errors)
		Expect(err.Error()).To(ContainSubstring("Validation error: "))
		Expect(err.Error()).To(ContainSubstring("field1 must be present, "))
		Expect(err.Error()).To(ContainSubstring("field1 must be blue or chartreuse, "))
		Expect(err.Error()).To(ContainSubstring("field2 must be unhinged, "))
	})
})
