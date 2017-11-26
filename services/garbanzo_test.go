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

	It("fetches garbanzos by octo name", func() {
		var garbanzos []data.Garbanzo
		mockGarbanzoStore.FetchByOctoNameOutput.Garbanzos <- garbanzos
		err := errors.New("some error")
		mockGarbanzoStore.FetchByOctoNameOutput.Err <- err

		actualGarbanzos, actualErr := service.FetchByOctoName(ctx, "my-octo")

		Expect(actualGarbanzos).To(Equal(garbanzos))
		Expect(actualErr).To(Equal(err))

		Expect(mockGarbanzoStore.FetchByOctoNameCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockGarbanzoStore.FetchByOctoNameInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockGarbanzoStore.FetchByOctoNameInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualOctoName string
		Expect(mockGarbanzoStore.FetchByOctoNameInput.OctoName).To(Receive(&actualOctoName))
		Expect(actualOctoName).To(Equal("my-octo"))
	})

	It("fetches a garbanzo by API UUID", func() {
		garbanzo := data.Garbanzo{
			GarbanzoType: data.DESI,
		}
		mockGarbanzoStore.FetchByAPIUUIDAndOctoNameOutput.Garbanzo <- garbanzo
		err := errors.New("some error")
		mockGarbanzoStore.FetchByAPIUUIDAndOctoNameOutput.Err <- err
		apiUUID := uuid.NewV4()

		actualGarbanzo, actualErr := service.FetchByAPIUUIDAndOctoName(ctx, apiUUID, "my-octo")

		Expect(actualGarbanzo).To(Equal(garbanzo))
		Expect(actualErr).To(Equal(err))

		Expect(mockGarbanzoStore.FetchByAPIUUIDAndOctoNameCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockGarbanzoStore.FetchByAPIUUIDAndOctoNameInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockGarbanzoStore.FetchByAPIUUIDAndOctoNameInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualAPIUUID uuid.UUID
		Expect(mockGarbanzoStore.FetchByAPIUUIDAndOctoNameInput.ApiUUID).To(Receive(&actualAPIUUID))
		Expect(actualAPIUUID).To(Equal(apiUUID))
		var actualOctoName string
		Expect(mockGarbanzoStore.FetchByAPIUUIDAndOctoNameInput.OctoName).To(Receive(&actualOctoName))
		Expect(actualOctoName).To(Equal("my-octo"))
	})

	Describe("Create", func() {
		It("creates a garbanzo", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			octoId := 77
			mockOctoStore.FetchByNameOutput.Octo <- data.Octo{
				Id: octoId,
			}
			mockOctoStore.FetchByNameOutput.Err <- nil

			garbanzoId := 42
			mockGarbanzoStore.CreateOutput.GarbanzoId <- garbanzoId
			mockGarbanzoStore.CreateOutput.Err <- nil

			mockTx.CommitOutput.Err <- nil

			garbanzo := data.Garbanzo{
				GarbanzoType: data.DESI,
				DiameterMM:   0.1,
			}
			octoName := "kraken"
			actualGarbanzo, actualErr := service.Create(ctx, octoName, garbanzo)

			Expect(actualErr).To(BeNil())
			Expect(actualGarbanzo.Id).To(Equal(garbanzoId))
			Expect(actualGarbanzo.APIUUID).NotTo(Equal(uuid.UUID{}))
			Expect(actualGarbanzo.GarbanzoType).To(Equal(data.DESI))

			Expect(mockDB.BeginTxCalled).To(HaveLen(1))

			Expect(mockOctoStore.FetchByNameCalled).To(HaveLen(1))
			var actualCtx context.Context
			Expect(mockOctoStore.FetchByNameInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			var actualDB persistence.Database
			Expect(mockOctoStore.FetchByNameInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			var actualOctoName string
			Expect(mockOctoStore.FetchByNameInput.Name).To(Receive(&actualOctoName))
			Expect(actualOctoName).To(Equal(octoName))
			var actualSelectForUpdate bool
			Expect(mockOctoStore.FetchByNameInput.SelectForUpdate).To(Receive(&actualSelectForUpdate))
			Expect(actualSelectForUpdate).To(Equal(true))

			Expect(mockGarbanzoStore.CreateCalled).To(HaveLen(1))
			Expect(mockGarbanzoStore.CreateInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			Expect(mockGarbanzoStore.CreateInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			var persistedGarbanzo data.Garbanzo
			Expect(mockGarbanzoStore.CreateInput.Garbanzo).To(Receive(&persistedGarbanzo))
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

			_, err := service.Create(ctx, "kraken", garbanzo)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("don't bother"))
		})

		It("returns an error if it can't find the octo", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			mockOctoStore.FetchByNameOutput.Octo <- data.Octo{}
			mockOctoStore.FetchByNameOutput.Err <- errors.New("don't bother")

			mockTx.RollbackOutput.Err <- nil

			garbanzo := data.Garbanzo{
				GarbanzoType: data.DESI,
				DiameterMM:   0.1,
			}

			_, err := service.Create(ctx, "kraken", garbanzo)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("don't bother"))

			Expect(mockTx.RollbackCalled).To(HaveLen(1))
		})

		It("returns a validation error for an empty request", func() {
			garbanzo := data.Garbanzo{}

			_, err := service.Create(ctx, "kraken", garbanzo)
			Expect(err).To(HaveOccurred())
			validationErr, ok := err.(services.ValidationError)
			Expect(ok).To(BeTrue())
			errors := validationErr.Errors()
			Expect(errors).To(HaveLen(2))
			Expect(errors).To(Equal(map[string][]string{
				"GarbanzoType": {"must be present", "must be either 'DESI' or 'KABULI'"},
				"DiameterMM":   {"must be present", "must be a positive decimal value"},
			}))
		})

		It("returns a validation error for invalid values", func() {
			garbanzo := data.Garbanzo{
				GarbanzoType: 1,
				DiameterMM:   -1.2,
			}

			_, err := service.Create(ctx, "kraken", garbanzo)
			Expect(err).To(HaveOccurred())
			validationErr, ok := err.(services.ValidationError)
			Expect(ok).To(BeTrue())
			errors := validationErr.Errors()
			Expect(errors).To(HaveLen(2))
			Expect(errors).To(Equal(map[string][]string{
				"GarbanzoType": {"must be either 'DESI' or 'KABULI'"},
				"DiameterMM":   {"must be a positive decimal value"},
			}))
		})
	})

	It("deletes a garbanzo by API UUID", func() {
		err := errors.New("some error")
		mockGarbanzoStore.DeleteByAPIUUIDAndOctoNameOutput.Err <- err
		apiUUID := uuid.NewV4()

		actualErr := service.DeleteByAPIUUIDAndOctoName(ctx, apiUUID, "my-octo")

		Expect(actualErr).To(Equal(err))

		Expect(mockGarbanzoStore.DeleteByAPIUUIDAndOctoNameCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockGarbanzoStore.DeleteByAPIUUIDAndOctoNameInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockGarbanzoStore.DeleteByAPIUUIDAndOctoNameInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualAPIUUID uuid.UUID
		Expect(mockGarbanzoStore.DeleteByAPIUUIDAndOctoNameInput.ApiUUID).To(Receive(&actualAPIUUID))
		Expect(actualAPIUUID).To(Equal(apiUUID))
		var actualOctoName string
		Expect(mockGarbanzoStore.DeleteByAPIUUIDAndOctoNameInput.OctoName).To(Receive(&actualOctoName))
		Expect(actualOctoName).To(Equal("my-octo"))
	})
})
