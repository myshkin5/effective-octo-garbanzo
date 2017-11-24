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
		database    persistence.Database
		store       persistence.GarbanzoStore
		octoId      int
		otherOctoId int
	)

	BeforeEach(func() {
		var err error
		database, err = persistence.Open()
		Expect(err).NotTo(HaveOccurred())

		Expect(cleanDatabase(database)).To(Succeed())

		octoId, err = persistence.OctoStore{}.Create(ctx, database, data.Octo{
			Name: "test-octo",
		})
		Expect(err).NotTo(HaveOccurred())

		otherOctoId, err = persistence.OctoStore{}.Create(ctx, database, data.Octo{
			Name: "test-octo-2",
		})
		Expect(err).NotTo(HaveOccurred())

		store = persistence.GarbanzoStore{}
	})

	Describe("FetchAll", func() {
		It("fetches no garbanzos when there are none", func() {
			garbanzos, err := store.FetchAll(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(0))
		})

		It("fetches all the garbanzos", func() {
			apiUUID1 := uuid.NewV4()
			garbanzo1 := data.Garbanzo{
				APIUUID:      apiUUID1,
				GarbanzoType: data.DESI,
				OctoId:       octoId,
				DiameterMM:   4.2,
			}
			apiUUID2 := uuid.NewV4()
			garbanzo2 := data.Garbanzo{
				APIUUID:      apiUUID2,
				GarbanzoType: data.KABULI,
				OctoId:       octoId,
				DiameterMM:   6.4,
			}

			garbanzoId1, err := store.Create(ctx, database, garbanzo1)
			Expect(err).NotTo(HaveOccurred())
			garbanzoId2, err := store.Create(ctx, database, garbanzo2)
			Expect(err).NotTo(HaveOccurred())

			garbanzos, err := store.FetchAll(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(2))

			Expect(garbanzos[0].Id).To(Equal(garbanzoId1))
			Expect(garbanzos[0].APIUUID).To(Equal(apiUUID1))
			Expect(garbanzos[0].GarbanzoType).To(Equal(data.DESI))
			Expect(garbanzos[0].OctoId).To(Equal(octoId))
			Expect(garbanzos[0].DiameterMM).To(BeNumerically("~", 4.2, 0.000001))

			Expect(garbanzos[1].Id).To(Equal(garbanzoId2))
			Expect(garbanzos[1].APIUUID).To(Equal(apiUUID2))
			Expect(garbanzos[1].GarbanzoType).To(Equal(data.KABULI))
			Expect(garbanzos[1].OctoId).To(Equal(octoId))
			Expect(garbanzos[1].DiameterMM).To(BeNumerically("~", 6.4, 0.000001))
		})
	})

	Describe("FetchByAPIUUID", func() {
		It("returns not found when fetching an unknown garbanzo", func() {
			_, err := store.FetchByAPIUUID(ctx, database, uuid.NewV4())

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := data.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: data.DESI,
				OctoId:       octoId,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.Create(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.GarbanzoType).To(Equal(data.DESI))
			Expect(fetchedGarbanzo.OctoId).To(Equal(octoId))
			Expect(fetchedGarbanzo.DiameterMM).To(BeNumerically("~", 4.2, 0.000001))
		})
	})

	Describe("Create", func() {
		It("creates a new garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := data.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: data.DESI,
				OctoId:       octoId,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.Create(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchByAPIUUID(ctx, database, apiUUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.GarbanzoType).To(Equal(data.DESI))
			Expect(fetchedGarbanzo.OctoId).To(Equal(octoId))
			Expect(fetchedGarbanzo.DiameterMM).To(BeNumerically("~", 4.2, 0.000001))
		})

		It("ignores the supplied id on create", func() {
			ignoredId := 82475928
			garbanzo := data.Garbanzo{
				Id:           ignoredId,
				APIUUID:      uuid.NewV4(),
				GarbanzoType: data.DESI,
				OctoId:       octoId,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.Create(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())
			Expect(garbanzoId).NotTo(Equal(ignoredId))
		})
	})

	Describe("DeleteByAPIUUID", func() {
		It("returns not found when deleting an unknown garbanzo", func() {
			err := store.DeleteByAPIUUID(ctx, database, uuid.NewV4())

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := data.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: data.DESI,
				OctoId:       octoId,
				DiameterMM:   4.2,
			}

			_, err := store.Create(ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			Expect(store.DeleteByAPIUUID(ctx, database, apiUUID)).To(Succeed())

			err = store.DeleteByAPIUUID(ctx, database, apiUUID)
			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})

	Describe("DeleteByOctoId", func() {
		It("returns no error when attempting to delete garbanzos from an unknown octo", func() {
			Expect(store.DeleteByOctoId(ctx, database, 2488)).To(Succeed())
		})

		It("returns no error when attempting to delete garbanzos from an octo with no garbanzos", func() {
			Expect(store.DeleteByOctoId(ctx, database, octoId)).To(Succeed())
		})

		It("deletes some garbanzos", func() {
			apiUUID1 := uuid.NewV4()
			garbanzo1 := data.Garbanzo{
				APIUUID:      apiUUID1,
				GarbanzoType: data.DESI,
				OctoId:       octoId,
				DiameterMM:   4.2,
			}
			apiUUID2 := uuid.NewV4()
			garbanzo2 := data.Garbanzo{
				APIUUID:      apiUUID2,
				GarbanzoType: data.KABULI,
				OctoId:       octoId,
				DiameterMM:   3.8,
			}
			apiUUID3 := uuid.NewV4()
			garbanzo3 := data.Garbanzo{
				APIUUID:      apiUUID3,
				GarbanzoType: data.DESI,
				OctoId:       otherOctoId,
				DiameterMM:   5.6,
			}

			_, err := store.Create(ctx, database, garbanzo1)
			Expect(err).NotTo(HaveOccurred())
			_, err = store.Create(ctx, database, garbanzo2)
			Expect(err).NotTo(HaveOccurred())
			octoId3, err := store.Create(ctx, database, garbanzo3)
			Expect(err).NotTo(HaveOccurred())
			garbanzo3.Id = octoId3

			Expect(store.DeleteByOctoId(ctx, database, octoId)).To(Succeed())

			garbanzos, err := store.FetchAll(ctx, database)
			Expect(err).NotTo(HaveOccurred())
			Expect(garbanzos).To(Equal([]data.Garbanzo{garbanzo3}))
		})
	})
})
