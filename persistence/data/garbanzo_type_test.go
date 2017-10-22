package data_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

var _ = Describe("GarbanzoType", func() {
	Describe("GarbanzoTypeFromString", func() {
		Context("happy path", func() {
			It("returns the DESI type", func() {
				gType, err := data.GarbanzoTypeFromString("DESI")
				Expect(err).ToNot(HaveOccurred())
				Expect(gType).To(Equal(data.DESI))
			})

			It("returns the KABULI type", func() {
				gType, err := data.GarbanzoTypeFromString("KABULI")
				Expect(err).ToNot(HaveOccurred())
				Expect(gType).To(Equal(data.KABULI))
			})
		})

		It("unhappy path", func() {
			_, err := data.GarbanzoTypeFromString("bogus")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid garbanzo type: bogus"))
		})
	})

	Describe("String", func() {
		It("returns the expected string", func() {
			Expect(data.DESI.String()).To(Equal("DESI"))
			Expect(data.KABULI.String()).To(Equal("KABULI"))
			Expect(data.GarbanzoType(42).String()).To(Equal("GarbanzoType(42)"))
		})
	})
})
