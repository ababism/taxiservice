package postgre

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"music-snap/pkg/app"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"testing"
	"time"
)

func initializeUserPortRepo() (ports.UserRepository, func() error, func() error, error) {
	db, closeDB, cleanDB, err := setupTestDB()
	if err != nil {
		return userRepository{}, nil, nil, app.NewError(0, "can't initialize DB", "initializeDB failed to initialize DB", err)
	}

	repo := newUserRepository(db)

	return repo, closeDB, cleanDB, nil
}

func initializeUserRepo() (userRepository, func() error, func() error, error) {
	db, closeDB, cleanDB, err := setupTestDB()
	if err != nil {
		return userRepository{}, nil, nil, app.NewError(0, "can't initialize DB", "initializeDB failed to initialize DB", err)
	}

	repo := newUserRepository(db)

	return repo, closeDB, cleanDB, nil
}

func TestUserRepository(t *testing.T) {
	repo, closeDB, cleanDB, err := initializeUserRepo()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer closeDB()
	defer cleanDB()

	profile := domain.Profile{
		ID:            uuid.New(),
		Nickname:      "testuser",
		AvatarURL:     "http://example.com/avatar.png",
		BackgroundURL: "http://example.com/bg.png",
		Bio:           "Test bio",
	}
	testUser := domain.User{
		Profile:      profile,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Roles:        domain.NewRoles([]string{domain.UserRole, domain.AdminRole}),
	}
	// CREATE ---------------------------------------------------------------------------
	t.Run("Test user create", func(t *testing.T) {
		t.Run("Test user create good", func(t *testing.T) {
			ctx := context.Background()

			createdUser, err := repo.Create(ctx, testUser)

			require.NoError(t, err)
			assert.Equal(t, testUser.ID, createdUser.ID)
			assert.Equal(t, testUser.Nickname, createdUser.Nickname)
			assert.Equal(t, testUser.Email, createdUser.Email)
			assert.Equal(t, testUser.AvatarURL, createdUser.AvatarURL)
			assert.Equal(t, testUser.BackgroundURL, createdUser.BackgroundURL)
			assert.Equal(t, testUser.Bio, createdUser.Bio)
			assert.True(t, createdUser.CreatedAt != time.Time{})
			assert.True(t, createdUser.UpdatedAt != time.Time{})

			testUser = createdUser
			fmt.Printf("createdUser: %v\n", createdUser)

		})

		t.Run("Test with duplicate email error", func(t *testing.T) {
			ctx := context.Background()
			profile := domain.Profile{
				ID:            uuid.New(),
				Nickname:      "duplicateuser",
				AvatarURL:     "http://example.com/avatar.png",
				BackgroundURL: "http://example.com/bg.png",
				Bio:           "Test bio",
			}
			user := domain.User{
				Profile:      profile,
				Email:        "duplicateemail@example.com",
				PasswordHash: "hashedpassword",
				Roles:        domain.NewRoles([]string{"role1", "role2"}),
			}

			_, _ = repo.Create(ctx, user) // Insert the testUser once

			_, err := repo.Create(ctx, user) // Attempt to insert the same testUser again

			assert.Error(t, err)
		})

		t.Run("Test empty bio", func(t *testing.T) {
			ctx := context.Background()

			profile := domain.Profile{
				ID:            uuid.New(),
				Nickname:      "emtybiouser",
				AvatarURL:     "http://example.com/avatar.png",
				BackgroundURL: "http://example.com/bg.png",
				Bio:           "",
			}

			user := domain.User{
				Profile:      profile,
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
				Roles:        domain.NewRoles([]string{}), // Empty roles
			}

			_, err := repo.Create(ctx, user)

			assert.Error(t, err)
		})
	})
	// GET ---------------------------------------------------------------------------
	t.Run("Test user get", func(t *testing.T) {
		t.Run("Test user get by ID", func(t *testing.T) {
			ctx := context.Background()

			retrievedUser, err := repo.GetByID(ctx, testUser.ID)
			require.NoError(t, err)
			assert.Equal(t, testUser.ID, retrievedUser.ID)
			assert.Equal(t, testUser.Nickname, retrievedUser.Nickname)
			assert.Equal(t, testUser.Email, retrievedUser.Email)
			assert.Equal(t, testUser.AvatarURL, retrievedUser.AvatarURL)
			assert.Equal(t, testUser.BackgroundURL, retrievedUser.BackgroundURL)
			assert.Equal(t, testUser.Bio, retrievedUser.Bio)
			assert.Equal(t, testUser.CreatedAt, retrievedUser.CreatedAt)
			assert.Equal(t, testUser.UpdatedAt, retrievedUser.UpdatedAt)
			assert.Equal(t, testUser.Roles, retrievedUser.Roles)
		},
		)
	})

	t.Run("Test user get", func(t *testing.T) {
		t.Run("Test user get by ID", func(t *testing.T) {
			ctx := context.Background()

			resProfile, err := repo.GetProfile(ctx, testUser.ID)
			require.NoError(t, err)
			assert.Equal(t, testUser.ID, resProfile.ID)
			assert.Equal(t, testUser.Nickname, resProfile.Nickname)

			assert.Equal(t, testUser.AvatarURL, resProfile.AvatarURL)
			assert.Equal(t, testUser.BackgroundURL, resProfile.BackgroundURL)
			assert.Equal(t, testUser.Bio, resProfile.Bio)
			assert.Equal(t, testUser.CreatedAt, resProfile.CreatedAt)
			assert.Equal(t, testUser.UpdatedAt, resProfile.UpdatedAt)
		},
		)
	})

	t.Run("Test user get by filter", func(t *testing.T) {
		t.Run("Test user get by filter ID", func(t *testing.T) {
			ctx := context.Background()

			f := domain.UserFilter{
				ID:        &testUser.ID,
				Nickname:  nil,
				Email:     nil,
				Limit:     nil,
				LastUUID:  nil,
				CreatedAt: nil,
				UpdatedAt: nil,
			}
			retrievedUser, err := repo.GetByParam(ctx, f)
			require.NoError(t, err)
			assert.Equal(t, testUser.ID, retrievedUser.ID)
			assert.Equal(t, testUser.Nickname, retrievedUser.Nickname)
			assert.Equal(t, testUser.Email, retrievedUser.Email)
			assert.Equal(t, testUser.AvatarURL, retrievedUser.AvatarURL)
			assert.Equal(t, testUser.BackgroundURL, retrievedUser.BackgroundURL)
			assert.Equal(t, testUser.Bio, retrievedUser.Bio)
			assert.Equal(t, testUser.CreatedAt, retrievedUser.CreatedAt)
			assert.Equal(t, testUser.UpdatedAt, retrievedUser.UpdatedAt)
			assert.Equal(t, testUser.Roles, retrievedUser.Roles)
		},
		)

		t.Run("Test user get by filter Nickname", func(t *testing.T) {
			ctx := context.Background()

			f := domain.UserFilter{
				ID:       nil,
				Nickname: &testUser.Nickname,
				Email:    nil,
			}
			retrievedUser, err := repo.GetByParam(ctx, f)
			require.NoError(t, err)
			assert.Equal(t, testUser.ID, retrievedUser.ID)
			assert.Equal(t, testUser.Nickname, retrievedUser.Nickname)
			assert.Equal(t, testUser.Email, retrievedUser.Email)
			assert.Equal(t, testUser.AvatarURL, retrievedUser.AvatarURL)
			assert.Equal(t, testUser.BackgroundURL, retrievedUser.BackgroundURL)
			assert.Equal(t, testUser.Bio, retrievedUser.Bio)
			assert.Equal(t, testUser.CreatedAt, retrievedUser.CreatedAt)
			assert.Equal(t, testUser.UpdatedAt, retrievedUser.UpdatedAt)
			assert.Equal(t, testUser.Roles, retrievedUser.Roles)
		},
		)

		t.Run("Test user get by filter Email", func(t *testing.T) {
			ctx := context.Background()

			f := domain.UserFilter{
				ID:       nil,
				Nickname: nil,
				Email:    &testUser.Email,
			}
			retrievedUser, err := repo.GetByParam(ctx, f)
			require.NoError(t, err)
			assert.Equal(t, testUser.ID, retrievedUser.ID)
			assert.Equal(t, testUser.Nickname, retrievedUser.Nickname)
			assert.Equal(t, testUser.Email, retrievedUser.Email)
			assert.Equal(t, testUser.AvatarURL, retrievedUser.AvatarURL)
			assert.Equal(t, testUser.BackgroundURL, retrievedUser.BackgroundURL)
			assert.Equal(t, testUser.Bio, retrievedUser.Bio)
			assert.Equal(t, testUser.CreatedAt, retrievedUser.CreatedAt)
			assert.Equal(t, testUser.UpdatedAt, retrievedUser.UpdatedAt)
			assert.Equal(t, testUser.Roles, retrievedUser.Roles)
		},
		)
	})

	// UPDATE ---------------------------------------------------------------------------
	t.Run("Test user update", func(t *testing.T) {
		t.Run("Test user update profile good", func(t *testing.T) {
			ctx := context.Background()

			newUser := testUser.Copy()
			newUser.Nickname = "updateduser"
			newUser.Bio = "Updated bio"

			updatedUser, err := repo.UpdateProfile(ctx, newUser)

			require.NoError(t, err)
			assert.Equal(t, newUser.ID, updatedUser.ID)
			assert.Equal(t, newUser.Nickname, updatedUser.Nickname)
			assert.Equal(t, newUser.Email, updatedUser.Email)
			assert.Equal(t, newUser.AvatarURL, updatedUser.AvatarURL)
			assert.Equal(t, newUser.BackgroundURL, updatedUser.BackgroundURL)
			assert.Equal(t, newUser.Bio, updatedUser.Bio)
			assert.Equal(t, newUser.CreatedAt, updatedUser.CreatedAt)
			assert.True(t, newUser.UpdatedAt != time.Time{})
			assert.True(t, updatedUser.UpdatedAt != testUser.UpdatedAt)

		})

		t.Run("Test user update user good", func(t *testing.T) {
			ctx := context.Background()

			newUser := testUser.Copy()
			newUser.Email = "updatedtest@example.com"
			newUser.Nickname = "updateduser"
			newUser.Bio = "Updated bio"

			updatedUser, err := repo.UpdateUser(ctx, newUser)

			require.NoError(t, err)
			assert.Equal(t, newUser.ID, updatedUser.ID)
			assert.Equal(t, newUser.Nickname, updatedUser.Nickname)
			assert.Equal(t, newUser.Email, updatedUser.Email)
			assert.Equal(t, newUser.AvatarURL, updatedUser.AvatarURL)
			assert.Equal(t, newUser.BackgroundURL, updatedUser.BackgroundURL)
			assert.Equal(t, newUser.Bio, updatedUser.Bio)
			assert.Equal(t, newUser.CreatedAt, updatedUser.CreatedAt)
			assert.True(t, newUser.UpdatedAt != time.Time{})
			assert.True(t, updatedUser.UpdatedAt != testUser.UpdatedAt)

		})
	})

	// DELETE ---------------------------------------------------------------------------
	t.Run("Test user delete", func(t *testing.T) {
		t.Run("Test user delete good", func(t *testing.T) {
			ctx := context.Background()

			err := repo.Delete(ctx, testUser.ID)
			require.NoError(t, err)

			// Try to get the deleted user
			_, err = repo.GetByID(ctx, testUser.ID)
			assert.Error(t, err)
		})
	})

	// GETLIST ---------------------------------------------------------------------------
	t.Run("Test user get list", func(t *testing.T) {
		ctx := context.Background()
		// Create a few test users
		for i := 0; i < 5; i++ {
			userForSearch := domain.User{
				Profile: domain.Profile{
					ID:            uuid.New(),
					Nickname:      fmt.Sprintf("specialtestuser%d", i),
					AvatarURL:     "http://example.com/avatar.png",
					BackgroundURL: "http://example.com/bg.png",
					Bio:           "Test bio",
				},
				Email:        fmt.Sprintf("test@example.com%d", i),
				PasswordHash: "hashedpassword",
				Roles:        domain.NewRoles([]string{domain.UserRole, domain.AdminRole}),
			}
			_, err := repo.Create(ctx, userForSearch)
			require.NoError(t, err)
		}

		t.Run("Test user get list good", func(t *testing.T) {
			ctx := context.Background()

			pag := domain.UUIDPagination{
				Limit:    10,
				LastUUID: uuid.Nil,
			}

			users, newLastUUID, err := repo.GetList(ctx, "special", pag)
			require.NoError(t, err)
			assert.NotEmpty(t, users)
			assert.True(t, len(users) == 5)
			assert.NotEqual(t, users[1].Nickname, uuid.Nil)
			assert.Equal(t, newLastUUID, users[len(users)-1].ID)

			//pag = domain.UUIDPagination{
			//	Limit:    2,
			//	LastUUID: uuid.Nil,
			//}
			//
			//users, newLastUUID, err = repo.GetList(ctx, "special", pag)
			//require.NoError(t, err)
			//assert.NotEmpty(t, users)
			//assert.True(t, len(users) == 2)
			//assert.Equal(t, newLastUUID, uuid.Nil)

		})

		t.Run("Test user get list pagination iteration", func(t *testing.T) {
			ctx := context.Background()

			pag := domain.UUIDPagination{
				Limit:    2,
				LastUUID: uuid.Nil,
			}

			var allUsers []domain.User
			for {
				users, newLastUUID, err := repo.GetList(ctx, "special", pag)
				require.NoError(t, err)
				if len(users) == 0 {
					break
				}
				allUsers = append(allUsers, users...)
				pag.LastUUID = newLastUUID
			}

			assert.Equal(t, len(allUsers), 5)
		})
	})
	// SUBS ---------------------------------------------------------------------------
	t.Run("Test user subscriptions", func(t *testing.T) {
		// Create two test users for subscription testing
		ctx := context.Background()

		var (
			testSub1 domain.Subscription
			testSub2 domain.Subscription
		)

		subscriber := domain.User{
			Profile: domain.Profile{
				ID:            uuid.New(),
				Nickname:      "subscriber",
				AvatarURL:     "http://example.com/avatar.png",
				BackgroundURL: "http://example.com/bg.png",
				Bio:           "Test bio",
			},
			Email:        "subscriber@example.com",
			PasswordHash: "hashedpassword",
			Roles:        domain.NewRoles([]string{domain.UserRole}),
		}

		followed := domain.User{
			Profile: domain.Profile{
				ID:            uuid.New(),
				Nickname:      "followed",
				AvatarURL:     "http://example.com/avatar.png",
				BackgroundURL: "http://example.com/bg.png",
				Bio:           "Test bio",
			},
			Email:        "followed@example.com",
			PasswordHash: "hashedpassword",
			Roles:        domain.NewRoles([]string{domain.UserRole}),
		}

		subscriber, err := repo.Create(ctx, subscriber)
		require.NoError(t, err)
		followed, err = repo.Create(ctx, followed)
		require.NoError(t, err)

		t.Run("Test create subscription", func(t *testing.T) {
			sub := domain.Subscription{
				SubscriberID: subscriber.ID,
				FollowedID:   followed.ID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			createdSub, err := repo.CreateSub(ctx, sub)
			testSub1 = createdSub

			require.NoError(t, err)
			assert.Equal(t, subscriber.ID, createdSub.SubscriberID)
			assert.Equal(t, followed.ID, createdSub.FollowedID)
			assert.True(t, createdSub.CreatedAt != time.Time{})
			assert.True(t, createdSub.UpdatedAt != time.Time{})
		})

		t.Run("Test get subscription", func(t *testing.T) {
			sub, err := repo.GetSub(ctx, subscriber.ID, followed.ID)
			require.NoError(t, err)
			assert.Equal(t, subscriber.ID, sub.SubscriberID)
			assert.Equal(t, followed.ID, sub.FollowedID)
		})

		t.Run("Test list subscriptions", func(t *testing.T) {
			// Create another user to follow
			anotherFollowed := domain.User{
				Profile: domain.Profile{
					ID:            uuid.New(),
					Nickname:      "anotherfollowed",
					AvatarURL:     "http://example.com/avatar.png",
					BackgroundURL: "http://example.com/bg.png",
					Bio:           "Test bio",
				},
				Email:        "anotherfollowed@example.com",
				PasswordHash: "hashedpassword",
				Roles:        domain.NewRoles([]string{domain.UserRole}),
			}
			anotherFollowed, err := repo.Create(ctx, anotherFollowed)
			require.NoError(t, err)

			// Create another subscription
			sub := domain.Subscription{
				SubscriberID: subscriber.ID,
				FollowedID:   anotherFollowed.ID,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			createdSub2, err := repo.CreateSub(ctx, sub)
			testSub2 = createdSub2

			require.NoError(t, err)

			initialPag := domain.IDPagination{
				Limit:  10,
				LastID: 0,
			}
			// Test listing subscriptions
			subs, pag, err := repo.ListSubscriptions(ctx, subscriber.ID, uuid.Nil, initialPag)
			require.NoError(t, err)
			assert.Equal(t, 2, len(subs))
			assert.Equal(t, testSub2.ID, pag.LastID)
		})

		t.Run("Test delete subscription", func(t *testing.T) {

			_, err := repo.DeleteSub(ctx, testSub1)
			require.NoError(t, err)
			//assert.Equal(t, subscriber.ID, deletedSub.SubscriberID)
			//assert.Equal(t, followed.ID, deletedSub.FollowedID)

			// Verify subscription is deleted
			_, err = repo.GetSub(ctx, subscriber.ID, followed.ID)
			assert.Error(t, err)

		})
	})
}
