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
		mockOctoStore     *mockOctoStore
		mockGarbanzoStore *mockGarbanzoStore
		mockDB            *mockDatabase
		mockTx            *mockDatabase
		service           *services.GarbanzoService
		ctx               context.Context
	)

	BeforeEach(func() {
		mockOctoStore = newMockOctoStore()
		mockGarbanzoStore = newMockGarbanzoStore()
		mockDB = newMockDatabase()
		mockTx = newMockDatabase()
		ctx = context.Background()

		service = services.NewGarbanzoService(mockOctoStore, mockGarbanzoStore, mockDB)
	})

	It("fetches all garbanzos", func() {
		var garbanzos []data.Garbanzo
		mockGarbanzoStore.FetchAllGarbanzosOutput.Garbanzos <- garbanzos
		err := errors.New("some error")
		mockGarbanzoStore.FetchAllGarbanzosOutput.Err <- err

		actualGarbanzos, actualErr := service.FetchAllGarbanzos(ctx)

		Expect(actualGarbanzos).To(Equal(garbanzos))
		Expect(actualErr).To(Equal(err))

		Expect(mockGarbanzoStore.FetchAllGarbanzosCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockGarbanzoStore.FetchAllGarbanzosInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockGarbanzoStore.FetchAllGarbanzosInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
	})

	It("fetches a garbanzo by API UUID", func() {
		garbanzo := data.Garbanzo{
			GarbanzoType: data.DESI,
		}
		mockGarbanzoStore.FetchGarbanzoByAPIUUIDOutput.Garbanzo <- garbanzo
		err := errors.New("some error")
		mockGarbanzoStore.FetchGarbanzoByAPIUUIDOutput.Err <- err
		apiUUID := uuid.NewV4()

		actualGarbanzo, actualErr := service.FetchGarbanzoByAPIUUID(ctx, apiUUID)

		Expect(actualGarbanzo).To(Equal(garbanzo))
		Expect(actualErr).To(Equal(err))

		Expect(mockGarbanzoStore.FetchGarbanzoByAPIUUIDCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockGarbanzoStore.FetchGarbanzoByAPIUUIDInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockGarbanzoStore.FetchGarbanzoByAPIUUIDInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualAPIUUID uuid.UUID
		Expect(mockGarbanzoStore.FetchGarbanzoByAPIUUIDInput.ApiUUID).To(Receive(&actualAPIUUID))
		Expect(actualAPIUUID).To(Equal(apiUUID))
	})

	Describe("CreateGarbanzo", func() {
		It("creates a garbanzo", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			octoId := 77
			mockOctoStore.FetchOctoByNameOutput.Octo <- data.Octo{
				Id: octoId,
			}
			mockOctoStore.FetchOctoByNameOutput.Err <- nil

			garbanzoId := 42
			mockGarbanzoStore.CreateGarbanzoOutput.GarbanzoId <- garbanzoId
			mockGarbanzoStore.CreateGarbanzoOutput.Err <- nil

			mockTx.CommitOutput.Err <- nil

			garbanzo := data.Garbanzo{
				GarbanzoType: data.DESI,
				DiameterMM:   0.1,
			}
			octoName := "kraken"
			actualGarbanzo, actualErr := service.CreateGarbanzo(ctx, octoName, garbanzo)

			Expect(actualErr).To(BeNil())
			Expect(actualGarbanzo.Id).To(Equal(garbanzoId))
			Expect(actualGarbanzo.APIUUID).NotTo(Equal(uuid.UUID{}))
			Expect(actualGarbanzo.GarbanzoType).To(Equal(data.DESI))

			Expect(mockDB.BeginTxCalled).To(HaveLen(1))

			Expect(mockOctoStore.FetchOctoByNameCalled).To(HaveLen(1))
			var actualCtx context.Context
			Expect(mockOctoStore.FetchOctoByNameInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			var actualDB persistence.Database
			Expect(mockOctoStore.FetchOctoByNameInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			var actualOctoName string
			Expect(mockOctoStore.FetchOctoByNameInput.Name).To(Receive(&actualOctoName))
			Expect(actualOctoName).To(Equal(octoName))

			Expect(mockGarbanzoStore.CreateGarbanzoCalled).To(HaveLen(1))
			Expect(mockGarbanzoStore.CreateGarbanzoInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			Expect(mockGarbanzoStore.CreateGarbanzoInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			var persistedGarbanzo data.Garbanzo
			Expect(mockGarbanzoStore.CreateGarbanzoInput.Garbanzo).To(Receive(&persistedGarbanzo))
			Expect(persistedGarbanzo.APIUUID).To(Equal(actualGarbanzo.APIUUID))
			Expect(persistedGarbanzo.GarbanzoType).To(Equal(garbanzo.GarbanzoType))
			Expect(persistedGarbanzo.OctoId).To(Equal(octoId))

			Expect(mockTx.CommitCalled).To(HaveLen(1))
		})

		It("returns an error if it can't start a transaction", func() {
			mockDB.BeginTxOutput.Database <- nil
			mockDB.BeginTxOutput.Err <- errors.New("don't bother")
			garbanzo := data.Garbanzo{
				GarbanzoType: data.DESI,
				DiameterMM:   0.1,
			}

			_, err := service.CreateGarbanzo(ctx, "kraken", garbanzo)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("don't bother"))
		})

		It("returns an error if it can't find the octo", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			mockOctoStore.FetchOctoByNameOutput.Octo <- data.Octo{}
			mockOctoStore.FetchOctoByNameOutput.Err <- errors.New("don't bother")

			mockTx.RollbackOutput.Err <- nil

			garbanzo := data.Garbanzo{
				GarbanzoType: data.DESI,
				DiameterMM:   0.1,
			}

			_, err := service.CreateGarbanzo(ctx, "kraken", garbanzo)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("don't bother"))

			Expect(mockTx.RollbackCalled).To(HaveLen(1))
		})

		It("returns a validation error for an empty request", func() {
			garbanzo := data.Garbanzo{}

			_, err := service.CreateGarbanzo(ctx, "kraken", garbanzo)
			Expect(err).To(HaveOccurred())
			validationErr, ok := err.(services.ValidationError)
			Expect(ok).To(BeTrue())
			Expect(validationErr.Errors()).To(HaveLen(2))
			Expect(validationErr.Errors()[0]).To(Equal("'type' is required and must be either 'DESI' or 'KABULI'"))
			Expect(validationErr.Errors()[1]).To(Equal("'diameter-mm' is required to be a positive decimal value"))
		})
	})

	It("deletes a garbanzo by API UUID", func() {
		err := errors.New("some error")
		mockGarbanzoStore.DeleteGarbanzoByAPIUUIDOutput.Err <- err
		apiUUID := uuid.NewV4()

		actualErr := service.DeleteGarbanzoByAPIUUID(ctx, apiUUID)

		Expect(actualErr).To(Equal(err))

		Expect(mockGarbanzoStore.DeleteGarbanzoByAPIUUIDCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockGarbanzoStore.DeleteGarbanzoByAPIUUIDInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockGarbanzoStore.DeleteGarbanzoByAPIUUIDInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualAPIUUID uuid.UUID
		Expect(mockGarbanzoStore.DeleteGarbanzoByAPIUUIDInput.ApiUUID).To(Receive(&actualAPIUUID))
		Expect(actualAPIUUID).To(Equal(apiUUID))
	})
})
