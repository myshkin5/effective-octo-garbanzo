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
			garbanzo1 := persistence.Garbanzo{
				APIUUID:      apiUUID1,
				GarbanzoType: persistence.DESI,
				DiameterMM:   4.2,
			}
			apiUUID2 := uuid.NewV4()
			garbanzo2 := persistence.Garbanzo{
				APIUUID:      apiUUID2,
				GarbanzoType: persistence.KABULI,
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
			Expect(garbanzos[0].GarbanzoType).To(Equal(persistence.DESI))
			Expect(garbanzos[0].DiameterMM).To(BeNumerically("~", 4.2, 0.000001))

			Expect(garbanzos[1].Id).To(Equal(garbanzoId2))
			Expect(garbanzos[1].APIUUID).To(Equal(apiUUID2))
			Expect(garbanzos[1].GarbanzoType).To(Equal(persistence.KABULI))
			Expect(garbanzos[1].DiameterMM).To(BeNumerically("~", 6.4, 0.000001))
		})
	})

	Context("FetchGarbanzoByAPIUUID", func() {
		It("returns not found when fetching an unknown garbanzo", func() {
			_, err := store.FetchGarbanzoByAPIUUID(ctx, database, uuid.NewV4())

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := persistence.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: persistence.DESI,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchGarbanzoByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.GarbanzoType).To(Equal(persistence.DESI))
			Expect(fetchedGarbanzo.DiameterMM).To(BeNumerically("~", 4.2, 0.000001))
		})
	})

	Context("CreateGarbanzo", func() {
		It("creates a new garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := persistence.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: persistence.DESI,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.CreateGarbanzo(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchGarbanzoByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.GarbanzoType).To(Equal(persistence.DESI))
			Expect(fetchedGarbanzo.DiameterMM).To(BeNumerically("~", 4.2, 0.000001))
		})

		It("ignores the supplied id on create", func() {
			ignoredId := 82475928
			garbanzo := persistence.Garbanzo{
				Id:           ignoredId,
				APIUUID:      uuid.NewV4(),
				GarbanzoType: persistence.DESI,
				DiameterMM:   4.2,
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
				APIUUID:      apiUUID,
				GarbanzoType: persistence.DESI,
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
