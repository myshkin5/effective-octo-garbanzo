package services_test

//go:generate hel

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
	"github.com/myshkin5/effective-octo-garbanzo/services"
)

var _ = Describe("Octo", func() {
	var (
		mockOctoStore     *mockOctoStore
		mockGarbanzoStore *mockGarbanzoStore
		mockDB            *mockDatabase
		mockTx            *mockDatabase
		service           *services.OctoService
		ctx               context.Context
	)

	BeforeEach(func() {
		mockOctoStore = newMockOctoStore()
		mockGarbanzoStore = newMockGarbanzoStore()
		mockDB = newMockDatabase()
		mockTx = newMockDatabase()
		ctx = context.Background()

		service = services.NewOctoService(mockOctoStore, mockGarbanzoStore, mockDB)
	})

	It("fetches all octos", func() {
		var octos []data.Octo
		mockOctoStore.FetchAllOutput.Octos <- octos
		err := errors.New("some error")
		mockOctoStore.FetchAllOutput.Err <- err

		actualOctos, actualErr := service.FetchAll(ctx)

		Expect(actualOctos).To(Equal(octos))
		Expect(actualErr).To(Equal(err))

		Expect(mockOctoStore.FetchAllCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockOctoStore.FetchAllInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockOctoStore.FetchAllInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
	})

	It("fetches a octo by name", func() {
		octo := data.Octo{
			Name: "kraken",
		}
		mockOctoStore.FetchByNameOutput.Octo <- octo
		err := errors.New("some error")
		mockOctoStore.FetchByNameOutput.Err <- err

		actualOcto, actualErr := service.FetchByName(ctx, "kraken")

		Expect(actualOcto).To(Equal(octo))
		Expect(actualErr).To(Equal(err))

		Expect(mockOctoStore.FetchByNameCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockOctoStore.FetchByNameInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockOctoStore.FetchByNameInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualName string
		Expect(mockOctoStore.FetchByNameInput.Name).To(Receive(&actualName))
		Expect(actualName).To(Equal("kraken"))
	})

	Describe("Create", func() {
		It("creates a octo", func() {
			octoId := 42
			mockOctoStore.CreateOutput.OctoId <- octoId
			mockOctoStore.CreateOutput.Err <- nil
			octo := data.Octo{
				Name: "kraken",
			}

			actualOcto, actualErr := service.Create(ctx, octo)

			Expect(actualErr).NotTo(HaveOccurred())
			Expect(actualOcto.Id).To(Equal(octoId))
			Expect(actualOcto.Name).To(Equal("kraken"))

			Expect(mockOctoStore.CreateCalled).To(HaveLen(1))
			var actualDB persistence.Database
			Expect(mockOctoStore.CreateInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockDB))
			var actualCtx context.Context
			Expect(mockOctoStore.CreateInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			var persistedOcto data.Octo
			Expect(mockOctoStore.CreateInput.Octo).To(Receive(&persistedOcto))
			Expect(persistedOcto.Name).To(Equal(actualOcto.Name))
		})

		It("returns a validation error for an empty octo name", func() {
			octo := data.Octo{}

			_, err := service.Create(ctx, octo)
			Expect(err).To(HaveOccurred())
			validationErr, ok := err.(services.ValidationError)
			Expect(ok).To(BeTrue())
			errors := validationErr.Errors()
			Expect(errors).To(HaveLen(1))
			Expect(errors).To(Equal(map[string][]string{"Name": {
				"must be present",
				"must match regular expression '^[\\w-]+$'",
			}}))
		})

		It("returns a validation error for an octo name with invalid characters", func() {
			octo := data.Octo{
				Name: " 283",
			}

			_, err := service.Create(ctx, octo)
			Expect(err).To(HaveOccurred())
			validationErr, ok := err.(services.ValidationError)
			Expect(ok).To(BeTrue())
			errors := validationErr.Errors()
			Expect(errors).To(HaveLen(1))
			Expect(errors).To(Equal(map[string][]string{"Name": {"must match regular expression '^[\\w-]+$'"}}))
		})
	})

	Describe("DeleteByName", func() {
		It("returns an error if it can't start a transaction", func() {
			mockDB.BeginTxOutput.Database <- nil
			mockDB.BeginTxOutput.Err <- errors.New("don't bother")

			err := service.DeleteByName(ctx, "kraken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("don't bother"))
		})

		It("rolls back and returns an error if it can't delete the child garbanzos", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			mockGarbanzoStore.DeleteByOctoNameOutput.Err <- errors.New("don't bother")

			mockTx.RollbackOutput.Err <- nil

			err := service.DeleteByName(ctx, "kraken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("don't bother"))

			Expect(mockTx.RollbackCalled).To(HaveLen(1))
		})

		It("rolls back and returns an error if it can't delete the octo", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			mockGarbanzoStore.DeleteByOctoNameOutput.Err <- nil

			mockOctoStore.DeleteByNameOutput.Err <- errors.New("some error")

			mockTx.RollbackOutput.Err <- nil

			err := service.DeleteByName(ctx, "kraken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("some error"))

			Expect(mockTx.RollbackCalled).To(HaveLen(1))
		})

		It("deletes a octo by name", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			mockGarbanzoStore.DeleteByOctoNameOutput.Err <- nil

			mockOctoStore.DeleteByNameOutput.Err <- nil

			mockTx.CommitOutput.Err <- nil

			actualErr := service.DeleteByName(ctx, "kraken")
			Expect(actualErr).NotTo(HaveOccurred())

			Expect(mockGarbanzoStore.DeleteByOctoNameCalled).To(HaveLen(1))
			var actualDB persistence.Database
			Expect(mockGarbanzoStore.DeleteByOctoNameInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			var actualCtx context.Context
			Expect(mockGarbanzoStore.DeleteByOctoNameInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			var actualName string
			Expect(mockGarbanzoStore.DeleteByOctoNameInput.OctoName).To(Receive(&actualName))
			Expect(actualName).To(Equal("kraken"))

			Expect(mockOctoStore.DeleteByNameCalled).To(HaveLen(1))
			Expect(mockOctoStore.DeleteByNameInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			Expect(mockOctoStore.DeleteByNameInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			Expect(mockOctoStore.DeleteByNameInput.Name).To(Receive(&actualName))
			Expect(actualName).To(Equal("kraken"))

			Expect(mockTx.CommitCalled).To(HaveLen(1))
		})
	})
})
