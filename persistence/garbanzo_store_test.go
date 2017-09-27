package persistence_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/satori/go.uuid"
)

var _ = Describe("GarbanzoStore Integration", func() {
	ctx := context.Background()
	var (
		database persistence.Database
		store    persistence.GarbanzoStore
	)

	BeforeEach(func() {
		var err error
		database, err = persistence.Open()
		Expect(err).NotTo(HaveOccurred())

		query := "delete from garbanzo"
		_, err = database.ExecContext(ctx, query)
		Expect(err).NotTo(HaveOccurred())

		store = persistence.GarbanzoStore{}
	})

	Context("FetchAllGarbanzos", func() {
		It("fetches no garbanzos when there are none", func() {
			garbanzos, err := store.FetchAllGarbanzos(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(0))
		})

		It("fetches all the garbanzos", func() {
			apiUUID1 := uuid.NewV4()
			firstName1 := "Joe"
			lastName1 := "Schmoe"
			garbanzo1 := persistence.Garbanzo{
				APIUUID:   apiUUID1,
				FirstName: firstName1,
				LastName:  lastName1,
			}
			apiUUID2 := uuid.NewV4()
			firstName2 := "Marty"
			lastName2 := "Blarty"
			garbanzo2 := persistence.Garbanzo{
				APIUUID:   apiUUID2,
				FirstName: firstName2,
				LastName:  lastName2,
			}

			garbanzoId1, err := store.CreateGarbanzo(ctx, database, garbanzo1)
			Expect(err).NotTo(HaveOccurred())
			garbanzoId2, err := store.CreateGarbanzo(ctx, database, garbanzo2)
			Expect(err).NotTo(HaveOccurred())

			garbanzos, err := store.FetchAllGarbanzos(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(2))

			Expect(garbanzos[0].Id).To(Equal(garbanzoId1))
			Expect(garbanzos[0].APIUUID).To(Equal(apiUUID1))
			Expect(garbanzos[0].FirstName).To(Equal(firstName1))
			Expect(garbanzos[0].LastName).To(Equal(lastName1))

			Expect(garbanzos[1].Id).To(Equal(garbanzoId2))
			Expect(garbanzos[1].APIUUID).To(Equal(apiUUID2))
			Expect(garbanzos[1].FirstName).To(Equal(firstName2))
			Expect(garbanzos[1].LastName).To(Equal(lastName2))
		})
	})

	Context("FetchGarbanzoByAPIUUID", func() {
		It("returns not found when fetching an unknown garbanzo", func() {
			_, err := store.FetchGarbanzoByAPIUUID(ctx, database, uuid.NewV4())

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a garbanzo", func() {
			apiUUID := uuid.NewV4()
			firstName := "Joe"
			lastName := "Schmoe"
			garbanzo := persistence.Garbanzo{
				APIUUID:   apiUUID,
				FirstName: firstName,
				LastName:  lastName,
			}

			garbanzoId, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchGarbanzoByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.FirstName).To(Equal(firstName))
			Expect(fetchedGarbanzo.LastName).To(Equal(lastName))
		})
	})

	Context("CreateGarbanzo", func() {
		It("creates a new garbanzo", func() {
			apiUUID := uuid.NewV4()
			firstName := "Joe"
			lastName := "Schmoe"
			garbanzo := persistence.Garbanzo{
				APIUUID:   apiUUID,
				FirstName: firstName,
				LastName:  lastName,
			}

			garbanzoId, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchGarbanzoByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.FirstName).To(Equal(firstName))
			Expect(fetchedGarbanzo.LastName).To(Equal(lastName))
		})

		It("ignores the supplied id on create", func() {
			ignoredId := 82475928
			garbanzo := persistence.Garbanzo{
				Id:        ignoredId,
				APIUUID:   uuid.NewV4(),
				FirstName: "Joe",
				LastName:  "Schmoe",
			}

			garbanzoId, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())
			Expect(garbanzoId).NotTo(Equal(ignoredId))
		})
	})

	Context("DeleteGarbanzoByAPIUUID", func() {
		It("returns not found when deleting an unknown garbanzo", func() {
			err := store.DeleteGarbanzoByAPIUUID(ctx, database, uuid.NewV4())

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := persistence.Garbanzo{
				APIUUID:   apiUUID,
				FirstName: "Joe",
				LastName:  "Schmoe",
			}

			_, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			err = store.DeleteGarbanzoByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())

			err = store.DeleteGarbanzoByAPIUUID(ctx, database, apiUUID)
			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})
})
