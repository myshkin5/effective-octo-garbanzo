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

	Context("FetchAllUsers", func() {
		It("fetches no users when there are none", func() {
			users, err := persistence.FetchAllUsers(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(users).To(HaveLen(0))
		})

		It("fetches all the users", func() {
			firstName1 := "Joe"
			lastName1 := "Schmoe"
			user1 := persistence.User{
				FirstName: firstName1,
				LastName:  lastName1,
			}
			firstName2 := "Marty"
			lastName2 := "Blarty"
			user2 := persistence.User{
				FirstName: firstName2,
				LastName:  lastName2,
			}

			userId1, err := persistence.CreateUser(ctx, database, user1)
			Expect(err).NotTo(HaveOccurred())
			userId2, err := persistence.CreateUser(ctx, database, user2)
			Expect(err).NotTo(HaveOccurred())

			users, err := persistence.FetchAllUsers(ctx, database)
			Expect(err).NotTo(HaveOccurred())

			Expect(users).To(HaveLen(2))

			Expect(users[0].Id).To(Equal(userId1))
			Expect(users[0].FirstName).To(Equal(firstName1))
			Expect(users[0].LastName).To(Equal(lastName1))

			Expect(users[1].Id).To(Equal(userId2))
			Expect(users[1].FirstName).To(Equal(firstName2))
			Expect(users[1].LastName).To(Equal(lastName2))
		})
	})

	Context("FetchUserById", func() {
		It("returns not found when fetching an unknown user", func() {
			_, err := persistence.FetchUserById(ctx, database, 727385737)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("fetches a user", func() {
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
	})

	Context("CreateUser", func() {
		It("creates a new user", func() {
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
			user := persistence.User{
				Id:        ignoredId,
				FirstName: "Joe",
				LastName:  "Schmoe",
			}

			userId, err := persistence.CreateUser(ctx, database, user)
			Expect(err).NotTo(HaveOccurred())
			Expect(userId).NotTo(Equal(ignoredId))
		})
	})

	Context("DeleteUserById", func() {
		It("returns not found when deleting an unknown user", func() {
			err := persistence.DeleteUserById(ctx, database, 727385737)

			Expect(err).To(Equal(persistence.ErrNotFound))
		})

		It("deletes a user", func() {
			user := persistence.User{
				FirstName: "Joe",
				LastName:  "Schmoe",
			}

			userId, err := persistence.CreateUser(ctx, database, user)
			Expect(err).NotTo(HaveOccurred())

			err = persistence.DeleteUserById(ctx, database, userId)
			Expect(err).NotTo(HaveOccurred())

			err = persistence.DeleteUserById(ctx, database, userId)
			Expect(err).To(Equal(persistence.ErrNotFound))
		})
	})
})
