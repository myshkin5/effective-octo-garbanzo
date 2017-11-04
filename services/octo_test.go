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
		mockStore.FetchAllOctosOutput.Octos <- octos
		err := errors.New("some error")
		mockStore.FetchAllOctosOutput.Err <- err

		actualOctos, actualErr := service.FetchAllOctos(ctx)

		Expect(actualOctos).To(Equal(octos))
		Expect(actualErr).To(Equal(err))

		Expect(mockStore.FetchAllOctosCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.FetchAllOctosInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.FetchAllOctosInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
	})

	It("fetches a octo by name", func() {
		octo := data.Octo{
			Name: "kraken",
		}
		mockStore.FetchOctoByNameOutput.Octo <- octo
		err := errors.New("some error")
		mockStore.FetchOctoByNameOutput.Err <- err

		actualOcto, actualErr := service.FetchOctoByName(ctx, "kraken")

		Expect(actualOcto).To(Equal(octo))
		Expect(actualErr).To(Equal(err))

		Expect(mockStore.FetchOctoByNameCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.FetchOctoByNameInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.FetchOctoByNameInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualName string
		Expect(mockStore.FetchOctoByNameInput.Name).To(Receive(&actualName))
		Expect(actualName).To(Equal("kraken"))
	})

	Describe("CreateOcto", func() {
		It("creates a octo", func() {
			octoId := 42
			mockStore.CreateOctoOutput.OctoId <- octoId
			mockStore.CreateOctoOutput.Err <- nil
			octo := data.Octo{
				Name: "kraken",
			}

			actualOcto, actualErr := service.CreateOcto(ctx, octo)

			Expect(actualErr).NotTo(HaveOccurred())
			Expect(actualOcto.Id).To(Equal(octoId))
			Expect(actualOcto.Name).To(Equal("kraken"))

			Expect(mockStore.CreateOctoCalled).To(HaveLen(1))
			var actualDB persistence.Database
			Expect(mockStore.CreateOctoInput.Database).To(Receive(&actualDB))
			Expect(actualDB).To(Equal(mockDB))
			var actualCtx context.Context
			Expect(mockStore.CreateOctoInput.Ctx).To(Receive(&actualCtx))
			Expect(actualCtx).To(Equal(ctx))
			var persistedOcto data.Octo
			Expect(mockStore.CreateOctoInput.Octo).To(Receive(&persistedOcto))
			Expect(persistedOcto.Name).To(Equal(actualOcto.Name))
		})

		It("returns a validation error for an empty octo name", func() {
			octo := data.Octo{}

			_, err := service.CreateOcto(ctx, octo)
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

			_, err := service.CreateOcto(ctx, octo)
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
		mockStore.DeleteOctoByNameOutput.Err <- err

		actualErr := service.DeleteOctoByName(ctx, "kraken")

		Expect(actualErr).To(Equal(err))

		Expect(mockStore.DeleteOctoByNameCalled).To(HaveLen(1))
		var actualDB persistence.Database
		Expect(mockStore.DeleteOctoByNameInput.Database).To(Receive(&actualDB))
		Expect(actualDB).To(Equal(mockDB))
		var actualCtx context.Context
		Expect(mockStore.DeleteOctoByNameInput.Ctx).To(Receive(&actualCtx))
		Expect(actualCtx).To(Equal(ctx))
		var actualName string
		Expect(mockStore.DeleteOctoByNameInput.Name).To(Receive(&actualName))
		Expect(actualName).To(Equal("kraken"))
	})
})
