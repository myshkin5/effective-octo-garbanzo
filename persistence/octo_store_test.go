package persistence_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

var _ = Describe("OctoStore Integration", func() {
	ctx := context.Background()
	var (
		database persistence.Database
		store    persistence.OctoStore
	)

	BeforeEach(func() {
		var err error
		database, err = persistence.Open()
		Expect(err).NotTo(HaveOccurred())

		err = cleanDatabase(database)
		Expect(err).NotTo(HaveOccurred())

		store = persistence.OctoStore{}
	})

	Describe("FetchAllOctos", func() {
		It("fetches no octos when there are none", func() {
			octos, err := store.FetchAllOctos(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(octos).To(HaveLen(0))
		})

		It("fetches all the octos", func() {
			octo1 := data.Octo{
				Name: "kraken",
			}
			octo2 := data.Octo{
				Name: "cthulhu",
			}

			octoId1, err := store.CreateOcto(ctx, database, octo1)
			Expect(err).NotTo(HaveOccurred())
			octoId2, err := store.CreateOcto(ctx, database, octo2)
			Expect(err).NotTo(HaveOccurred())

			octos, err := store.FetchAllOctos(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(octos).To(HaveLen(2))

			Expect(octos[0].Id).To(Equal(octoId1))
			Expect(octos[0].Name).To(Equal("kraken"))

			Expect(octos[1].Id).To(Equal(octoId2))
			Expect(octos[1].Name).To(Equal("cthulhu"))
		})
	})

	Describe("FetchOctoByName", func() {
		It("returns not found when fetching an unknown octo", func() {
			_, err := store.FetchOctoByName(ctx, database, "squidward")

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a octo", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			octoId, err := store.CreateOcto(ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			fetchedOcto, err := store.FetchOctoByName(ctx, database, "kraken")
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedOcto.Id).To(Equal(octoId))
			Expect(fetchedOcto.Name).To(Equal("kraken"))
		})
	})

	Describe("CreateOcto", func() {
		It("creates a new octo", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			octoId, err := store.CreateOcto(ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			fetchedOcto, err := store.FetchOctoByName(ctx, database, "kraken")
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedOcto.Id).To(Equal(octoId))
			Expect(fetchedOcto.Name).To(Equal("kraken"))
		})

		It("ignores the supplied id on create", func() {
			ignoredId := 82475928
			octo := data.Octo{
				Id:   ignoredId,
				Name: "kraken",
			}

			octoId, err := store.CreateOcto(ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())
			Expect(octoId).NotTo(Equal(ignoredId))
		})
	})

	Describe("DeleteOctoByName", func() {
		It("returns not found when deleting an unknown octo", func() {
			err := store.DeleteOctoByName(ctx, database, "squidward")

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a octo", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			_, err := store.CreateOcto(ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			err = store.DeleteOctoByName(ctx, database, "kraken")
			Expect(err).NotTo(HaveOccurred())

			err = store.DeleteOctoByName(ctx, database, "kraken")
			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})
})
