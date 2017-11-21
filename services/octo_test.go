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
		mockStore *mockOctoStore
		mockDB    *mockDatabase
		service   *services.OctoService
		ctx       context.Context
	)

	BeforeEach(func() {
		mockStore = newMockOctoStore()
		mockDB = newMockDatabase()
		ctx = context.Background()

		service = services.NewOctoService(mockStore, mockDB)
	})

	It("fetches all octos", func() {
		var octos []data.Octo
		mockStore.FetchAllOutput.Octos <- octos
		err := errors.New("some error")
		mockStore.FetchAllOutput.Err <- err

		actualOctos, actualErr := service.FetchAll(ctx)

		Expect(actualOctos).To(Equal(octos))
		Expect(actualErr).To(Equal(err))

		Expect(mockStore.FetchAllCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.FetchAllInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.FetchAllInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
	})

	It("fetches a octo by name", func() {
		octo := data.Octo{
			Name: "kraken",
		}
		mockStore.FetchByNameOutput.Octo <- octo
		err := errors.New("some error")
		mockStore.FetchByNameOutput.Err <- err

		actualOcto, actualErr := service.FetchByName(ctx, "kraken")

		Expect(actualOcto).To(Equal(octo))
		Expect(actualErr).To(Equal(err))

		Expect(mockStore.FetchByNameCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.FetchByNameInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.FetchByNameInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualName string
		Expect(mockStore.FetchByNameInput.Name).To(Receive(&actualName))
		Expect(actualName).To(Equal("kraken"))
	})

	Describe("Create", func() {
		It("creates a octo", func() {
			octoId := 42
			mockStore.CreateOutput.OctoId <- octoId
			mockStore.CreateOutput.Err <- nil
			octo := data.Octo{
				Name: "kraken",
			}

			actualOcto, actualErr := service.Create(ctx, octo)

			Expect(actualErr).NotTo(HaveOccurred())
			Expect(actualOcto.Id).To(Equal(octoId))
			Expect(actualOcto.Name).To(Equal("kraken"))

			Expect(mockStore.CreateCalled).To(HaveLen(1))
			var actualDB persistence.Database
			Expect(mockStore.CreateInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockDB))
			var actualCtx context.Context
			Expect(mockStore.CreateInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			var persistedOcto data.Octo
			Expect(mockStore.CreateInput.Octo).To(Receive(&persistedOcto))
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

	It("deletes a octo by name", func() {
		err := errors.New("some error")
		mockStore.DeleteByNameOutput.Err <- err

		actualErr := service.DeleteByName(ctx, "kraken")

		Expect(actualErr).To(Equal(err))

		Expect(mockStore.DeleteByNameCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.DeleteByNameInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.DeleteByNameInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualName string
		Expect(mockStore.DeleteByNameInput.Name).To(Receive(&actualName))
		Expect(actualName).To(Equal("kraken"))
	})
})
