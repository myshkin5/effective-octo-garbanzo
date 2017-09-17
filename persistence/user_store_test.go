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
		database *sql.DB
	)

	BeforeEach(func() {
		var err error
		database, err = sql.Open(
			"postgres",
			"postgres://garbanzo:garbanzo-secret@localhost:5678/garbanzo?sslmode=disable")
		Expect(err).NotTo(HaveOccurred())

		query := "delete from garbanzo_user"
		_, err = database.ExecContext(ctx, query)
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns not found when fetching an unknown user", func() {
		_, err := persistence.FetchUserById(ctx, database, 727385737)

		Expect(err).To(Equal(persistence.ErrNotFound))
	})

	It("returns not found when deleting an unknown user", func() {
		err := persistence.DeleteUserById(ctx, database, 727385737)

		Expect(err).To(Equal(persistence.ErrNotFound))
	})

	It("creates and fetches a new user", func() {
		firstName := "Joe"
		lastName := "Schmoe"
		user := persistence.User{
			FirstName: firstName,
			LastName:  lastName,
		}

		userId, err := persistence.CreateUser(ctx, database, user)
		Expect(err).NotTo(HaveOccurred())

		fetchedUser, err := persistence.FetchUserById(ctx, database, userId)
		Expect(err).NotTo(HaveOccurred())
		Expect(fetchedUser.Id).To(Equal(userId))
		Expect(fetchedUser.FirstName).To(Equal(firstName))
		Expect(fetchedUser.LastName).To(Equal(lastName))
	})

	It("ignores the supplied id on create", func() {
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
	})

	It("creates and deletes a new user", func() {
		firstName := "Joe"
		lastName := "Schmoe"
		user := persistence.User{
			FirstName: firstName,
			LastName:  lastName,
		}

		userId, err := persistence.CreateUser(ctx, database, user)
		Expect(err).NotTo(HaveOccurred())

		err = persistence.DeleteUserById(ctx, database, userId)
		Expect(err).NotTo(HaveOccurred())

		err = persistence.DeleteUserById(ctx, database, userId)
		Expect(err).To(Equal(persistence.ErrNotFound))
	})
})
