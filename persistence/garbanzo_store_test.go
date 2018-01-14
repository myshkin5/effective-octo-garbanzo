package persistence_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

var _ = Describe("GarbanzoStore Integration", func() {
	var (
		database                                                   persistence.Database
		store                                                      persistence.GarbanzoStore
		org1Ctx, org2Ctx                                           context.Context
		org1Octo1, org1Octo2, org2Octo1                            data.Octo
		org1Octo1Garbanzo1, org1Octo1Garbanzo2, org2Octo1Garbanzo1 data.Garbanzo
	)

	BeforeEach(func() {
		var err error
		database, err = persistence.Open()
		Expect(err).NotTo(HaveOccurred())

		cleanDatabase(database)

		_, orgName := createOrg("garbanzo_store", database)
		_, orgName2 := createOrg("garbanzo_store2", database)

		org1Ctx = context.WithValue(ctx, persistence.OrgContextKey, orgName)
		org2Ctx = context.WithValue(ctx, persistence.OrgContextKey, orgName2)

		octoStore := persistence.OctoStore{}

		org1Octo1 = data.Octo{
			Name: "org1Octo1",
		}
		id, err := octoStore.Create(org1Ctx, database, org1Octo1)
		Expect(err).NotTo(HaveOccurred())
		org1Octo1.Id = id

		org1Octo2 = data.Octo{
			Name: "org1Octo2",
		}
		id, err = octoStore.Create(org1Ctx, database, org1Octo2)
		Expect(err).NotTo(HaveOccurred())
		org1Octo2.Id = id

		org2Octo1 = data.Octo{
			Name: "org2Octo1",
		}
		id, err = octoStore.Create(org2Ctx, database, org2Octo1)
		Expect(err).NotTo(HaveOccurred())
		org2Octo1.Id = id

		store = persistence.GarbanzoStore{}

		org1Octo1Garbanzo1 = data.Garbanzo{
			APIUUID:      uuid.NewV4(),
			GarbanzoType: data.DESI,
			OctoId:       org1Octo1.Id,
			DiameterMM:   4.2,
		}
		id, err = store.Create(org1Ctx, database, org1Octo1Garbanzo1)
		Expect(err).NotTo(HaveOccurred())
		org1Octo1Garbanzo1.Id = id

		org1Octo1Garbanzo2 = data.Garbanzo{
			APIUUID:      uuid.NewV4(),
			GarbanzoType: data.KABULI,
			OctoId:       org1Octo1.Id,
			DiameterMM:   6.4,
		}
		id, err = store.Create(org1Ctx, database, org1Octo1Garbanzo2)
		Expect(err).NotTo(HaveOccurred())
		org1Octo1Garbanzo2.Id = id

		org2Octo1Garbanzo1 = data.Garbanzo{
			APIUUID:      uuid.NewV4(),
			GarbanzoType: data.KABULI,
			OctoId:       org2Octo1.Id,
			DiameterMM:   6.4,
		}
		id, err = store.Create(org2Ctx, database, org2Octo1Garbanzo1)
		Expect(err).NotTo(HaveOccurred())
		org2Octo1Garbanzo1.Id = id
	})

	Describe("FetchByOctoName", func() {
		It("fetches no garbanzos when there are none", func() {
			garbanzos, err := store.FetchByOctoName(org1Ctx, database, org1Octo2.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(0))
		})

		It("fetches all the garbanzos", func() {
			garbanzos, err := store.FetchByOctoName(org1Ctx, database, org1Octo1.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(HaveLen(2))

			Expect(garbanzos[0].Id).To(Equal(org1Octo1Garbanzo1.Id))
			Expect(garbanzos[0].APIUUID).To(Equal(org1Octo1Garbanzo1.APIUUID))
			Expect(garbanzos[0].GarbanzoType).To(Equal(data.DESI))
			Expect(garbanzos[0].OctoId).To(Equal(org1Octo1.Id))
			Expect(garbanzos[0].DiameterMM).To(BeNumerically("~", 4.2, 0.000001))

			Expect(garbanzos[1].Id).To(Equal(org1Octo1Garbanzo2.Id))
			Expect(garbanzos[1].APIUUID).To(Equal(org1Octo1Garbanzo2.APIUUID))
			Expect(garbanzos[1].GarbanzoType).To(Equal(data.KABULI))
			Expect(garbanzos[1].OctoId).To(Equal(org1Octo1.Id))
			Expect(garbanzos[1].DiameterMM).To(BeNumerically("~", 6.4, 0.000001))
		})

		It("does not find garbanzos for another org", func() {
			garbanzos, err := store.FetchByOctoName(org2Ctx, database, org1Octo1.Name)
			Expect(err).NotTo(HaveOccurred())

			Expect(garbanzos).To(BeEmpty())
		})
	})

	Describe("FetchByAPIUUIDAndOctoName", func() {
		It("returns not found when fetching an unknown garbanzo", func() {
			_, err := store.FetchByAPIUUIDAndOctoName(org1Ctx, database, uuid.NewV4(), org1Octo1.Name)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("returns not found when fetching a garbanzo with the wrong octo name", func() {
			_, err := store.FetchByAPIUUIDAndOctoName(org1Ctx, database, org1Octo1Garbanzo1.APIUUID, org1Octo2.Name)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a garbanzo", func() {
			fetchedGarbanzo, err := store.FetchByAPIUUIDAndOctoName(org1Ctx, database, org1Octo1Garbanzo1.APIUUID, org1Octo1.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(org1Octo1Garbanzo1.Id))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(org1Octo1Garbanzo1.APIUUID))
			Expect(fetchedGarbanzo.GarbanzoType).To(Equal(data.DESI))
			Expect(fetchedGarbanzo.OctoId).To(Equal(org1Octo1.Id))
			Expect(fetchedGarbanzo.DiameterMM).To(BeNumerically("~", 4.2, 0.000001))
		})

		It("does not find garbanzos for another org", func() {
			_, err := store.FetchByAPIUUIDAndOctoName(org2Ctx, database, org1Octo1Garbanzo1.APIUUID, org1Octo1.Name)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})

	Describe("Create", func() {
		It("creates a new garbanzo", func() {
			apiUUID := uuid.NewV4()
			garbanzo := data.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: data.DESI,
				OctoId:       org1Octo1.Id,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.Create(org1Ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())

			fetchedGarbanzo, err := store.FetchByAPIUUIDAndOctoName(org1Ctx, database, apiUUID, org1Octo1.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedGarbanzo.Id).To(Equal(garbanzoId))
			Expect(fetchedGarbanzo.APIUUID).To(Equal(apiUUID))
			Expect(fetchedGarbanzo.GarbanzoType).To(Equal(data.DESI))
			Expect(fetchedGarbanzo.OctoId).To(Equal(org1Octo1.Id))
			Expect(fetchedGarbanzo.DiameterMM).To(BeNumerically("~", 4.2, 0.000001))
		})

		It("ignores the supplied id on create", func() {
			ignoredId := 82475928
			garbanzo := data.Garbanzo{
				Id:           ignoredId,
				APIUUID:      uuid.NewV4(),
				GarbanzoType: data.DESI,
				OctoId:       org1Octo2.Id,
				DiameterMM:   4.2,
			}

			garbanzoId, err := store.Create(org1Ctx, database, garbanzo)
			Expect(err).NotTo(HaveOccurred())
			Expect(garbanzoId).NotTo(Equal(ignoredId))

			garbanzos, err := store.FetchByOctoName(org1Ctx, database, org1Octo2.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(garbanzos)).To(Equal(1))
			Expect(garbanzos[0].Id).To(Equal(garbanzoId))
		})

		It("fails to create a garbanzo with a parent octo from another org", func() {
			apiUUID := uuid.NewV4()
			garbanzo := data.Garbanzo{
				APIUUID:      apiUUID,
				GarbanzoType: data.DESI,
				OctoId:       org1Octo1.Id,
				DiameterMM:   4.2,
			}

			_, err := store.Create(org2Ctx, database, garbanzo)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DeleteByAPIUUIDAndOctoName", func() {
		It("returns not found when deleting an unknown garbanzo", func() {
			err := store.DeleteByAPIUUIDAndOctoName(org1Ctx, database, uuid.NewV4(), org1Octo1.Name)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("returns not found when deleting a garbanzo with the wrong octo name", func() {
			err := store.DeleteByAPIUUIDAndOctoName(org1Ctx, database, org1Octo1Garbanzo1.APIUUID, org1Octo2.Name)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a garbanzo", func() {
			Expect(store.DeleteByAPIUUIDAndOctoName(org1Ctx, database, org1Octo1Garbanzo1.APIUUID, org1Octo1.Name)).To(Succeed())

			err := store.DeleteByAPIUUIDAndOctoName(org1Ctx, database, org1Octo1Garbanzo1.APIUUID, org1Octo1.Name)
			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("returns not found when deleting a garbanzo with the wrong org", func() {
			err := store.DeleteByAPIUUIDAndOctoName(org2Ctx, database, org1Octo1Garbanzo1.APIUUID, org1Octo1.Name)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})

	Describe("DeleteByOctoId", func() {
		It("returns no error when attempting to delete garbanzos from an unknown octo", func() {
			Expect(store.DeleteByOctoId(org1Ctx, database, 24889404)).To(Succeed())
		})

		It("returns no error when attempting to delete garbanzos from an octo with no garbanzos", func() {
			Expect(store.DeleteByOctoId(org1Ctx, database, org1Octo1.Id)).To(Succeed())
		})

		It("deletes some garbanzos", func() {
			org1Octo2Garbanzo1 := data.Garbanzo{
				APIUUID:      uuid.NewV4(),
				GarbanzoType: data.DESI,
				OctoId:       org1Octo2.Id,
				DiameterMM:   5.6,
			}

			id, err := store.Create(org1Ctx, database, org1Octo2Garbanzo1)
			Expect(err).NotTo(HaveOccurred())
			org1Octo2Garbanzo1.Id = id

			Expect(store.DeleteByOctoId(org1Ctx, database, org1Octo1.Id)).To(Succeed())

			garbanzos, err := store.FetchByOctoName(org1Ctx, database, org1Octo1.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(garbanzos)).To(Equal(0))

			garbanzos, err = store.FetchByOctoName(org1Ctx, database, org1Octo2.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(garbanzos).To(Equal([]data.Garbanzo{org1Octo2Garbanzo1}))
		})

		It("returns no error when deleting garbanzos with the wrong org (but doesn't actually delete anything)", func() {
			Expect(store.DeleteByOctoId(org2Ctx, database, org1Octo1.Id)).To(Succeed())

			garbanzos, err := store.FetchByOctoName(org1Ctx, database, org1Octo1.Name)
			Expect(err).NotTo(HaveOccurred())
			Expect(garbanzos).To(HaveLen(2))
		})
	})
})
