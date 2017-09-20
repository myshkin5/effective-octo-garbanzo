package persistence_test

import (
	"context"
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

var _ = Describe("GarbanzoStore Integration", func() {
	ctx := context.Background()
	var (
		database *sql.DB
	)

	BeforeEach(func() {
		var err error
		database, err = sql.Open(
			"postgres",
			"postgres://garbanzo:garbanzo-secret@localhost:5678/garbanzo?sslmode=disable")
		Expect(err).NotTo(HaveOccurred())

		query := "delete from garbanzo"
		_, err = database.ExecContext(ctx, query)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("FetchAllGarbanzos", func() {
		It("fetches no garbanzos when there are none", func() {
			garbanzos, err := persistence.FetchAllGarbanzos(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(0))
		})

		It("fetches all the garbanzos", func() {
			firstName1 := "Joe"
			lastName1 := "Schmoe"
			garbanzo1 := persistence.Garbanzo{
				FirstName: firstName1,
				LastName:  lastName1,
			}
			firstName2 := "Marty"
			lastName2 := "Blarty"
			garbanzo2 := persistence.Garbanzo{
				FirstName: firstName2,
				LastName:  lastName2,
			}

			garbanzoId1, err := persistence.CreateGarbanzo(ctx, database, garbanzo1)
			Expect(err).NotTo(HaveOccurred())
			garbanzoId2, err := persistence.CreateGarbanzo(ctx, database, garbanzo2)
			Expect(err).NotTo(HaveOccurred())

			garbanzos, err := persistence.FetchAllGarbanzos(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(2))

			Expect(garbanzos[0].Id).To(Equal(garbanzoId1))
			Expect(garbanzos[0].FirstName).To(Equal(firstName1))
			Expect(garbanzos[0].LastName).To(Equal(lastName1))

			Expect(garbanzos[1].Id).To(Equal(garbanzoId2))
			Expect(garbanzos[1].FirstName).To(Equal(firstName2))
			Expect(garbanzos[1].LastName).To(Equal(lastName2))
		})
	})

	Context("FetchGarbanzoById", func() {
		It("returns not found when fetching an unknown garbanzo", func() {
			_, err := persistence.FetchGarbanzoById(ctx, database, 727385737)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a garbanzo", func() {
			firstName := "Joe"
			lastName := "Schmoe"
			garbanzo := persistence.Garbanzo{
				FirstName: firstName,
				LastName:  lastName,
			}

			garbanzoId, err := persistence.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := persistence.FetchGarbanzoById(ctx, database, garbanzoId)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.FirstName).To(Equal(firstName))
			Expect(fetchedGarbanzo.LastName).To(Equal(lastName))
		})
	})

	Context("CreateGarbanzo", func() {
		It("creates a new garbanzo", func() {
			firstName := "Joe"
			lastName := "Schmoe"
			garbanzo := persistence.Garbanzo{
				FirstName: firstName,
				LastName:  lastName,
			}

			garbanzoId, err := persistence.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := persistence.FetchGarbanzoById(ctx, database, garbanzoId)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.FirstName).To(Equal(firstName))
			Expect(fetchedGarbanzo.LastName).To(Equal(lastName))
		})

		It("ignores the supplied id on create", func() {
			ignoredId := 82475928
			garbanzo := persistence.Garbanzo{
				Id:        ignoredId,
				FirstName: "Joe",
				LastName:  "Schmoe",
			}

			garbanzoId, err := persistence.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())
			Expect(garbanzoId).NotTo(Equal(ignoredId))
		})
	})

	Context("DeleteGarbanzoById", func() {
		It("returns not found when deleting an unknown garbanzo", func() {
			err := persistence.DeleteGarbanzoById(ctx, database, 727385737)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a garbanzo", func() {
			garbanzo := persistence.Garbanzo{
				FirstName: "Joe",
				LastName:  "Schmoe",
			}

			garbanzoId, err := persistence.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			err = persistence.DeleteGarbanzoById(ctx, database, garbanzoId)
			Expect(err).NotTo(HaveOccurred())

			err = persistence.DeleteGarbanzoById(ctx, database, garbanzoId)
			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})
})
