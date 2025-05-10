package postgre

import (
	"context"
	"github.com/google/uuid"
	_ "github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"music-snap/services/musicsnap/internal/domain"
	"testing"
)

func TestReactionRepository(t *testing.T) {
	repo, closeDB, cleanDB, err := initializeRepository()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer closeDB()
	defer cleanDB()

	// Create a test user first since reviews need a user

	require.NoError(t, err)

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
		Roles:        domain.NewRoles([]string{domain.UserRole}),
	}

	createdUser, err := repo.user.Create(context.Background(), testUser)
	require.NoError(t, err)

	testUser = createdUser

	testReview := domain.Review{
		UserID:    createdUser.ID,
		PieceID:   uuid.New().String(),
		Rating:    9,
		Content:   "Great piece of music!",
		Moderated: false,
		Published: true,
		PhotoURL:  "https://example.com/photo.jpg",
	}

	var createdReview domain.Review
	createdReview, err = repo.review.Create(context.Background(), testReview)
	require.NoError(t, err)

	testReview = createdReview

	testReaction := domain.Reaction{
		UserID:   testUser.ID,
		ReviewID: testReview.ID,
		Type:     domain.LikeReaction,
	}

	// CREATE ---------------------------------------------------------------------------
	t.Run("Test reaction create", func(t *testing.T) {
		t.Run("Test reaction create good", func(t *testing.T) {
			ctx := context.Background()

			crReaction, err := repo.reaction.Create(ctx, testReaction)

			require.NoError(t, err)
			assert.Equal(t, testReaction.UserID, crReaction.UserID)
			assert.Equal(t, testReaction.ReviewID, crReaction.ReviewID)
			assert.Equal(t, testReaction.Type, crReaction.Type)
			assert.NotEqual(t, testReaction.ID, crReaction.ID) // ID should be generated

			testReaction = crReaction

			// Verify the reaction was created by getting it
			reaction, err := repo.reaction.GetByID(ctx, testReaction.ID)
			require.NoError(t, err)
			assert.Equal(t, testReaction.UserID, reaction.UserID)
			assert.Equal(t, testReaction.ReviewID, reaction.ReviewID)
			assert.Equal(t, testReaction.Type, reaction.Type)
		})

		t.Run("Test create duplicate reaction", func(t *testing.T) {
			ctx := context.Background()

			// Try to create the same reaction again
			_, err := repo.reaction.Create(ctx, testReaction)
			assert.Error(t, err) // Should fail due to unique constraint
		})
	})

	// GET ---------------------------------------------------------------------------
	t.Run("Test reaction get", func(t *testing.T) {
		t.Run("Test get reaction by review", func(t *testing.T) {
			ctx := context.Background()

			reactions, err := repo.reaction.GetFromReview(ctx, testReview.ID)
			require.NoError(t, err)
			require.Len(t, reactions, 1)
			assert.Equal(t, testReaction.UserID, reactions[0].UserID)
			assert.Equal(t, testReaction.ReviewID, reactions[0].ReviewID)
			assert.Equal(t, testReaction.Type, reactions[0].Type)
		})

		t.Run("Test get reaction by non-existent review", func(t *testing.T) {
			ctx := context.Background()

			reactions, err := repo.reaction.GetFromReview(ctx, 999999)
			require.NoError(t, err)
			assert.Empty(t, reactions)
		})
	})

	// UPDATE ---------------------------------------------------------------------------
	t.Run("Test reaction update", func(t *testing.T) {
		t.Run("Test update reaction type", func(t *testing.T) {
			ctx := context.Background()

			// Get the created reaction first
			obtReact, err := repo.reaction.GetByID(ctx, testReaction.ID)
			require.NoError(t, err)

			// Update the reaction type

			obtReact.Type = domain.DislikeReaction

			updReact, err := repo.reaction.Update(ctx, obtReact)
			require.NoError(t, err)
			assert.Equal(t, testReaction.ID, updReact.ID)
			assert.Equal(t, testReaction.UserID, updReact.UserID)
			assert.Equal(t, testReaction.ReviewID, updReact.ReviewID)
			assert.Equal(t, domain.DislikeReaction, updReact.Type)

			// Verify the update
			resReact, err := repo.reaction.GetByID(ctx, testReaction.ID)
			require.NoError(t, err)
			assert.Equal(t, domain.DislikeReaction, resReact.Type)
		})

		t.Run("Test update non-existent reaction", func(t *testing.T) {
			ctx := context.Background()

			nonExistentReaction := domain.Reaction{
				ID:       999999,
				UserID:   createdUser.ID,
				ReviewID: createdReview.ID,
				Type:     domain.LikeReaction,
			}

			_, err = repo.reaction.Update(ctx, nonExistentReaction)
			assert.Error(t, err)
		})
	})

	// DELETE ---------------------------------------------------------------------------
	t.Run("Test reaction delete", func(t *testing.T) {
		t.Run("Test delete existing reaction", func(t *testing.T) {
			ctx := context.Background()

			err = repo.reaction.Delete(ctx, testReaction.ID)
			require.NoError(t, err)

			reaction, err := repo.reaction.GetByID(ctx, testReaction.ID)
			require.Error(t, err)
			assert.Empty(t, reaction)
		})

		t.Run("Test delete non-existent reaction", func(t *testing.T) {
			ctx := context.Background()

			err := repo.reaction.Delete(ctx, 999999)
			assert.Error(t, err)
		})
	})

	// Additional test cases for multiple reactions
	t.Run("Test multiple reactions on same review", func(t *testing.T) {
		ctx := context.Background()

		// Create another user
		secondUser := domain.User{
			Profile: domain.Profile{
				ID:            uuid.New(),
				Nickname:      "seconduser",
				AvatarURL:     "http://example.com/avatar2.png",
				BackgroundURL: "http://example.com/bg2.png",
				Bio:           "Second test bio",
			},
			Email:        "test2@example.com",
			PasswordHash: "hashedpassword2",
			Roles:        domain.NewRoles([]string{domain.UserRole}),
		}

		createdSecondUser, err := repo.user.Create(ctx, secondUser)
		require.NoError(t, err)

		// Create a reaction from the second user
		testReaction2 := domain.Reaction{
			UserID:   createdSecondUser.ID,
			ReviewID: createdReview.ID,
			Type:     domain.LikeReaction,
		}

		testReaction, err = repo.reaction.Create(ctx, testReaction)
		require.NoError(t, err)

		testReaction2, err = repo.reaction.Create(ctx, testReaction2)
		require.NoError(t, err)

		// Verify both reactions exist
		reactions, err := repo.reaction.GetFromReview(ctx, testReview.ID)
		require.NoError(t, err)
		assert.Len(t, reactions, 2)

		// Clean up
		for _, reaction := range reactions {
			err = repo.reaction.Delete(ctx, reaction.ID)
			require.NoError(t, err)
		}
	})
}
