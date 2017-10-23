package persistence_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
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

	Describe("FetchAllGarbanzos", func() {
		It("fetches no garbanzos when there are none", func() {
			garbanzos, err := store.FetchAllGarbanzos(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(0))
		})

		It("fetches all the garbanzos", func() {
			apiUUID1 := uuid.NewV4()
			garbanzo1 := data.Garbanzo{
				APIUUID:      apiUUID1,
				GarbanzoType: data.DESI,
				DiameterMM:   4.2,
			}
			apiUUID2 := uuid.NewV4()
			garbanzo2 := data.Garbanzo{
				APIUUID:      apiUUID2,
				GarbanzoType: data.KABULI,
				DiameterMM:   6.4,
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
			Expect(garbanzos[0].GarbanzoType).To(Equal(data.DESI))
			Expect(garbanzos[0].DiameterMM).To(BeNumerically("~", 4.2, 0.000001))

			Expect(garbanzos[1].Id).To(Equal(garbanzoId2))
			Expect(garbanzos[1].APIUUID).To(Equal(apiUUID2))
			Expect(garbanzos[1].GarbanzoType).To(Equal(data.KABULI))
			Expect(garbanzos[1].DiameterMM).To(BeNumerically("~", 6.4, 0.000001))
		})
	})

	Describe("FetchGarbanzoByAPIUUID", func() {
		It("returns not found when fetching an unknown garbanzo", func() {
			_, err := store.FetchGarbanzoByAPIUUID(ctx, database, uuid.NewV4())

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := data.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: data.DESI,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchGarbanzoByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.GarbanzoType).To(Equal(data.DESI))
			Expect(fetchedGarbanzo.DiameterMM).To(BeNumerically("~", 4.2, 0.000001))
		})
	})

	Describe("CreateGarbanzo", func() {
		It("creates a new garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := data.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: data.DESI,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchGarbanzoByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.GarbanzoType).To(Equal(data.DESI))
			Expect(fetchedGarbanzo.DiameterMM).To(BeNumerically("~", 4.2, 0.000001))
		})

		It("ignores the supplied id on create", func() {
			ignoredId := 82475928
			garbanzo := data.Garbanzo{
				Id:           ignoredId,
				APIUUID:      uuid.NewV4(),
				GarbanzoType: data.DESI,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())
			Expect(garbanzoId).NotTo(Equal(ignoredId))
		})
	})

	Describe("DeleteGarbanzoByAPIUUID", func() {
		It("returns not found when deleting an unknown garbanzo", func() {
			err := store.DeleteGarbanzoByAPIUUID(ctx, database, uuid.NewV4())

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := data.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: data.DESI,
				DiameterMM:   4.2,
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
