package services_test

//go:generate hel

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
	"github.com/myshkin5/effective-octo-garbanzo/services"
)

var _ = Describe("Garbanzo", func() {
	var (
		mockStore *mockGarbanzoStore
		mockDB    *mockDatabase
		service   *services.GarbanzoService
	)

	BeforeEach(func() {
		mockStore = newMockGarbanzoStore()
		mockDB = newMockDatabase()

		service = services.NewGarbanzoService(mockStore, mockDB)
	})

	It("fetches all garbanzos", func() {
		var garbanzos []data.Garbanzo
		mockStore.FetchAllGarbanzosOutput.Garbanzos <- garbanzos
		err := errors.New("some error")
		mockStore.FetchAllGarbanzosOutput.Err <- err
		ctx := context.TODO()

		actualGarbanzos, actualErr := service.FetchAllGarbanzos(ctx)

		Expect(actualGarbanzos).To(Equal(garbanzos))
		Expect(actualErr).To(Equal(err))

		Expect(mockStore.FetchAllGarbanzosCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.FetchAllGarbanzosInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.FetchAllGarbanzosInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
	})

	It("fetches a garbanzo by API UUID", func() {
		garbanzo := data.Garbanzo{
			GarbanzoType: data.DESI,
		}
		mockStore.FetchGarbanzoByAPIUUIDOutput.Garbanzo <- garbanzo
		err := errors.New("some error")
		mockStore.FetchGarbanzoByAPIUUIDOutput.Err <- err
		ctx := context.TODO()
		apiUUID := uuid.NewV4()

		actualGarbanzo, actualErr := service.FetchGarbanzoByAPIUUID(ctx, apiUUID)

		Expect(actualGarbanzo).To(Equal(garbanzo))
		Expect(actualErr).To(Equal(err))

		Expect(mockStore.FetchGarbanzoByAPIUUIDCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.FetchGarbanzoByAPIUUIDInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.FetchGarbanzoByAPIUUIDInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualAPIUUID uuid.UUID
		Expect(mockStore.FetchGarbanzoByAPIUUIDInput.ApiUUID).To(Receive(&actualAPIUUID))
		Expect(actualAPIUUID).To(Equal(apiUUID))
	})

	It("creates a garbanzo", func() {
		garbanzoId := 42
		mockStore.CreateGarbanzoOutput.GarbanzoId <- garbanzoId
		mockStore.CreateGarbanzoOutput.Err <- nil
		ctx := context.TODO()
		garbanzo := data.Garbanzo{
			GarbanzoType: data.DESI,
		}

		actualGarbanzo, actualErr := service.CreateGarbanzo(ctx, garbanzo)

		Expect(actualGarbanzo.Id).To(Equal(garbanzoId))
		Expect(actualErr).To(BeNil())
		Expect(actualGarbanzo.APIUUID).NotTo(Equal(uuid.UUID{}))
		Expect(actualGarbanzo.GarbanzoType).To(Equal(data.DESI))

		Expect(mockStore.CreateGarbanzoCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.CreateGarbanzoInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.CreateGarbanzoInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var persistedGarbanzo data.Garbanzo
		Expect(mockStore.CreateGarbanzoInput.Garbanzo).To(Receive(&persistedGarbanzo))
		Expect(persistedGarbanzo.APIUUID).To(Equal(actualGarbanzo.APIUUID))
		Expect(persistedGarbanzo.GarbanzoType).To(Equal(garbanzo.GarbanzoType))
	})

	It("deletes a garbanzo by API UUID", func() {
		err := errors.New("some error")
		mockStore.DeleteGarbanzoByAPIUUIDOutput.Err <- err
		ctx := context.TODO()
		apiUUID := uuid.NewV4()

		actualErr := service.DeleteGarbanzoByAPIUUID(ctx, apiUUID)

		Expect(actualErr).To(Equal(err))

		Expect(mockStore.DeleteGarbanzoByAPIUUIDCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.DeleteGarbanzoByAPIUUIDInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.DeleteGarbanzoByAPIUUIDInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualAPIUUID uuid.UUID
		Expect(mockStore.DeleteGarbanzoByAPIUUIDInput.ApiUUID).To(Receive(&actualAPIUUID))
		Expect(actualAPIUUID).To(Equal(apiUUID))
	})
})
