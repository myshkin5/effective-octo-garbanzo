package persistence_test

import (
	"context"
	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

var _ = Describe("UserStore Integration", func() {
	ctx := context.Background()
	var (
		database persistence.Database
	)

	BeforeEach(func() {
		var err error
		database, err = sql.Open("postgres", "postgres://localhost/garbanzo?sslmode=disable")
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns not found when fetching an unknown user", func() {
		// TODO: Eventually this will be run against a clean database to assure that this user never exists
		_, err := persistence.FetchUserById(ctx, database, 727385737)

		Expect(err).To(Equal(persistence.ErrNotFound))
	})

	It("returns not found when deleting an unknown user", func() {
		// TODO: Eventually this will be run against a clean database to assure that this user never exists
		err := persistence.DeleteUserById(ctx, database, 727385737)

		Expect(err).To(Equal(persistence.ErrNotFound))
	})

	It("creates, fetches, and deletes a new user", func() {
		// TODO: Break this up when we have a clean test database
		ignoredId := 82475928
		firstName := "Joe"
		lastName := "Schmoe"
		user := persistence.User{
			Id:        ignoredId,
			FirstName: firstName,
			LastName:  lastName,
		}

		userId, err := persistence.CreateUser(ctx, database, user)
		Expect(err).NotTo(HaveOccurred())
		Expect(userId).NotTo(Equal(ignoredId))

		fetchedUser, err := persistence.FetchUserById(ctx, database, userId)
		Expect(err).NotTo(HaveOccurred())
		Expect(fetchedUser.Id).To(Equal(userId))
		Expect(fetchedUser.FirstName).To(Equal(firstName))
		Expect(fetchedUser.LastName).To(Equal(lastName))

		err = persistence.DeleteUserById(ctx, database, userId)
		Expect(err).NotTo(HaveOccurred())

		err = persistence.DeleteUserById(ctx, database, userId)
		Expect(err).To(Equal(persistence.ErrNotFound))
	})
})
