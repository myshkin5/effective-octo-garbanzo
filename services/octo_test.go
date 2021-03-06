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

		It("rolls back and returns an error if it can't select the octo for update", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			mockOctoStore.FetchByNameOutput.Octo <- data.Octo{}
			mockOctoStore.FetchByNameOutput.Err <- errors.New("don't bother")

			mockTx.RollbackOutput.Err <- nil

			err := service.DeleteByName(ctx, "kraken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("don't bother"))

			Expect(mockTx.RollbackCalled).To(HaveLen(1))
		})

		It("rolls back and returns an error if it can't delete the child garbanzos", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			mockOctoStore.FetchByNameOutput.Octo <- data.Octo{
				Id:   282,
				Name: "kraken",
			}
			mockOctoStore.FetchByNameOutput.Err <- nil

			mockGarbanzoStore.DeleteByOctoIdOutput.Err <- errors.New("don't bother")

			mockTx.RollbackOutput.Err <- nil

			err := service.DeleteByName(ctx, "kraken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("don't bother"))

			Expect(mockTx.RollbackCalled).To(HaveLen(1))
		})

		It("rolls back and returns an error if it can't delete the octo", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			mockOctoStore.FetchByNameOutput.Octo <- data.Octo{
				Id:   282,
				Name: "kraken",
			}
			mockOctoStore.FetchByNameOutput.Err <- nil

			mockGarbanzoStore.DeleteByOctoIdOutput.Err <- nil

			mockOctoStore.DeleteByIdOutput.Err <- errors.New("some error")

			mockTx.RollbackOutput.Err <- nil

			err := service.DeleteByName(ctx, "kraken")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("some error"))

			Expect(mockTx.RollbackCalled).To(HaveLen(1))
		})

		It("deletes a octo by name", func() {
			mockDB.BeginTxOutput.Database <- mockTx
			mockDB.BeginTxOutput.Err <- nil

			id := 282
			mockOctoStore.FetchByNameOutput.Octo <- data.Octo{
				Id:   id,
				Name: "kraken",
			}
			mockOctoStore.FetchByNameOutput.Err <- nil

			mockGarbanzoStore.DeleteByOctoIdOutput.Err <- nil

			mockOctoStore.DeleteByIdOutput.Err <- nil

			mockTx.CommitOutput.Err <- nil

			actualErr := service.DeleteByName(ctx, "kraken")
			Expect(actualErr).NotTo(HaveOccurred())

			Expect(mockOctoStore.FetchByNameCalled).To(HaveLen(1))
			var actualDB persistence.Database
			Expect(mockOctoStore.FetchByNameInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			var actualCtx context.Context
			Expect(mockOctoStore.FetchByNameInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			var actualName string
			Expect(mockOctoStore.FetchByNameInput.Name).To(Receive(&actualName))
			Expect(actualName).To(Equal("kraken"))
			var actualSelectForUpdate bool
			Expect(mockOctoStore.FetchByNameInput.SelectForUpdate).To(Receive(&actualSelectForUpdate))
			Expect(actualSelectForUpdate).To(BeTrue())

			Expect(mockGarbanzoStore.DeleteByOctoIdCalled).To(HaveLen(1))
			Expect(mockGarbanzoStore.DeleteByOctoIdInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			Expect(mockGarbanzoStore.DeleteByOctoIdInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			var actualId int
			Expect(mockGarbanzoStore.DeleteByOctoIdInput.OctoId).To(Receive(&actualId))
			Expect(actualId).To(Equal(id))

			Expect(mockOctoStore.DeleteByIdCalled).To(HaveLen(1))
			Expect(mockOctoStore.DeleteByIdInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockTx))
			Expect(mockOctoStore.DeleteByIdInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			Expect(mockOctoStore.DeleteByIdInput.Id).To(Receive(&actualId))
			Expect(actualId).To(Equal(id))

			Expect(mockTx.CommitCalled).To(HaveLen(1))
		})
	})
})
