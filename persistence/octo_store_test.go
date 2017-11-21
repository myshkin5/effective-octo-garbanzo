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

		Expect(cleanDatabase(database)).To(Succeed())

		store = persistence.OctoStore{}
	})

	Describe("FetchAll", func() {
		It("fetches no octos when there are none", func() {
			octos, err := store.FetchAll(ctx, database)
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

			octoId1, err := store.Create(ctx, database, octo1)
			Expect(err).NotTo(HaveOccurred())
			octoId2, err := store.Create(ctx, database, octo2)
			Expect(err).NotTo(HaveOccurred())

			octos, err := store.FetchAll(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(octos).To(HaveLen(2))

			Expect(octos[0].Id).To(Equal(octoId1))
			Expect(octos[0].Name).To(Equal("kraken"))

			Expect(octos[1].Id).To(Equal(octoId2))
			Expect(octos[1].Name).To(Equal("cthulhu"))
		})
	})

	Describe("FetchByName", func() {
		It("returns not found when fetching an unknown octo", func() {
			_, err := store.FetchByName(ctx, database, "squidward")

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a octo", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			octoId, err := store.Create(ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			fetchedOcto, err := store.FetchByName(ctx, database, "kraken")
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedOcto.Id).To(Equal(octoId))
			Expect(fetchedOcto.Name).To(Equal("kraken"))
		})
	})

	Describe("Create", func() {
		It("creates a new octo", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			octoId, err := store.Create(ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			fetchedOcto, err := store.FetchByName(ctx, database, "kraken")
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

			octoId, err := store.Create(ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())
			Expect(octoId).NotTo(Equal(ignoredId))
		})
	})

	Describe("DeleteByName", func() {
		It("returns not found when deleting an unknown octo", func() {
			err := store.DeleteByName(ctx, database, "squidward")

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a octo", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			_, err := store.Create(ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			Expect(store.DeleteByName(ctx, database, "kraken")).To(Succeed())

			err = store.DeleteByName(ctx, database, "kraken")
			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})
})
