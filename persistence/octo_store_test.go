package persistence_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

var _ = Describe("OctoStore Integration", func() {
	var (
		database         persistence.Database
		store            persistence.OctoStore
		org1Ctx, org2Ctx context.Context
	)

	BeforeEach(func() {
		var err error
		database, err = persistence.Open()
		Expect(err).NotTo(HaveOccurred())

		cleanDatabase(database)

		_, orgName := createOrg("octo_store", database)
		_, orgName2 := createOrg("octo_store2", database)

		org1Ctx = context.WithValue(ctx, persistence.OrgContextKey, orgName)
		org2Ctx = context.WithValue(ctx, persistence.OrgContextKey, orgName2)

		store = persistence.OctoStore{}
	})

	Describe("FetchAll", func() {
		It("fetches no octos when there are none", func() {
			octos, err := store.FetchAll(org1Ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(octos).To(HaveLen(0))
		})

		It("fetches all the octos", func() {
			org1Octo1 := data.Octo{
				Name: "kraken",
			}
			org1Octo1Id, err := store.Create(org1Ctx, database, org1Octo1)
			Expect(err).NotTo(HaveOccurred())

			org1Octo2 := data.Octo{
				Name: "cthulhu",
			}
			org1Octo2Id, err := store.Create(org1Ctx, database, org1Octo2)
			Expect(err).NotTo(HaveOccurred())

			org2Octo1 := data.Octo{
				Name: "barry",
			}
			_, err = store.Create(org2Ctx, database, org2Octo1)
			Expect(err).NotTo(HaveOccurred())

			octos, err := store.FetchAll(org1Ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(octos).To(HaveLen(2))

			Expect(octos[0].Id).To(Equal(org1Octo1Id))
			Expect(octos[0].Name).To(Equal("kraken"))

			Expect(octos[1].Id).To(Equal(org1Octo2Id))
			Expect(octos[1].Name).To(Equal("cthulhu"))
		})
	})

	Describe("FetchByName", func() {
		Context("normal select", func() {
			It("returns not found when fetching an unknown octo", func() {
				_, err := store.FetchByName(org1Ctx, database, "squidward", false)

				Expect(err).To(Equal(persistence.ErrNotFound))
			})

			It("fetches a octo", func() {
				octo := data.Octo{
					Name: "kraken",
				}

				octoId, err := store.Create(org1Ctx, database, octo)
				Expect(err).NotTo(HaveOccurred())

				fetchedOcto, err := store.FetchByName(org1Ctx, database, "kraken", false)
				Expect(err).NotTo(HaveOccurred())
				Expect(fetchedOcto.Id).To(Equal(octoId))
				Expect(fetchedOcto.Name).To(Equal("kraken"))
			})

			It("does not find octos for another org", func() {
				octo := data.Octo{
					Name: "kraken",
				}

				_, err := store.Create(org1Ctx, database, octo)
				Expect(err).NotTo(HaveOccurred())

				_, err = store.FetchByName(org2Ctx, database, "kraken", false)
				Expect(err).To(Equal(persistence.ErrNotFound))
			})
		})

		Context("select for update", func() {
			It("locks records from other transactions", func() {
				octo := data.Octo{
					Name: "kraken",
				}

				_, err := store.Create(org1Ctx, database, octo)
				Expect(err).NotTo(HaveOccurred())

				tx1, err := database.BeginTx(ctx)
				Expect(err).NotTo(HaveOccurred())

				_, err = store.FetchByName(org1Ctx, tx1, "kraken", true)
				Expect(err).NotTo(HaveOccurred())

				tx2, err := database.BeginTx(ctx)
				Expect(err).NotTo(HaveOccurred())

				done := make(chan struct{}, 0)
				// NB: 0
				go func() {
					// NB: 1 -- everything written to between 1s and 2s must be unique to avoid data race errors
					_, err2 := persistence.OctoStore{}.FetchByName(org1Ctx, tx2, "kraken", true)
					Expect(err2).NotTo(HaveOccurred())
					close(done)
					// NB: 2
				}()

				// NB: 1
				Consistently(done).ShouldNot(BeClosed())

				err = tx1.Commit()
				Expect(err).NotTo(HaveOccurred())

				Eventually(done).Should(BeClosed())
				// NB: 2

				err = tx2.Commit()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Create", func() {
		It("creates a new octo", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			octoId, err := store.Create(org1Ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			fetchedOcto, err := store.FetchByName(org1Ctx, database, "kraken", false)
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

			octoId, err := store.Create(org1Ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())
			Expect(octoId).NotTo(Equal(ignoredId))
		})

		It("allows octos with the same name in different orgs", func() {
			octo1 := data.Octo{
				Name: "kraken",
			}

			_, err := store.Create(org1Ctx, database, octo1)
			Expect(err).NotTo(HaveOccurred())

			octo2 := data.Octo{
				Name: "kraken",
			}

			_, err = store.Create(org2Ctx, database, octo2)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("DeleteById", func() {
		It("returns not found when deleting an unknown octo", func() {
			err := store.DeleteById(org1Ctx, database, 82333455)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a octo", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			id, err := store.Create(org1Ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			Expect(store.DeleteById(org1Ctx, database, id)).To(Succeed())

			err = store.DeleteById(org1Ctx, database, id)
			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("returns not found when deleting an octo for another org", func() {
			octo := data.Octo{
				Name: "kraken",
			}

			id, err := store.Create(org1Ctx, database, octo)
			Expect(err).NotTo(HaveOccurred())

			err = store.DeleteById(org2Ctx, database, id)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})
})
