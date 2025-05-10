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
	"testing"
)

func initializeRepository() (repository, func() error, func() error, error) {
	db, closeDB, cleanDB, err := setupTestDB()
	if err != nil {
		return repository{}, nil, nil, app.NewError(0, "can't initialize DB", "initializeDB failed to initialize DB", err)
	}

	repo := newRepository(db)

	return repo, closeDB, cleanDB, nil
}

func TestReviewRepository(t *testing.T) {
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

	testReview := domain.Review{
		UserID:    createdUser.ID,
		PieceID:   uuid.New().String(),
		Rating:    9,
		Content:   "Great piece of music!",
		Moderated: false,
		Published: true,
		PhotoURL:  "https://example.com/photo.jpg",
	}

	testReview2 := domain.Review{}

	testReview3 := domain.Review{}

	// CREATE ---------------------------------------------------------------------------
	t.Run("Test review create", func(t *testing.T) {
		t.Run("Test review create good", func(t *testing.T) {
			ctx := context.Background()

			createdReview, err := repo.review.Create(ctx, testReview)

			require.NoError(t, err)
			assert.Equal(t, testReview.UserID, createdReview.UserID)
			assert.Equal(t, testReview.PieceID, createdReview.PieceID)
			assert.Equal(t, testReview.Rating, createdReview.Rating)
			assert.Equal(t, testReview.Content, createdReview.Content)
			assert.Equal(t, testReview.Moderated, createdReview.Moderated)
			assert.Equal(t, testReview.Published, createdReview.Published)
			assert.Equal(t, testReview.PhotoURL, createdReview.PhotoURL)
			assert.NotEqual(t, testReview.CreatedAt, createdReview.CreatedAt)
			assert.NotEqual(t, testReview.UpdatedAt, createdReview.UpdatedAt)

			testReview = createdReview
		})

		t.Run("Test with second review on same pieceID no error", func(t *testing.T) {
			ctx := context.Background()

			secondReview := testReview
			secondReview.ID = 2
			secondReview.Content = "Another great review!"

			createdReview, err := repo.review.Create(ctx, secondReview)
			require.NoError(t, err)
			assert.Equal(t, secondReview.UserID, createdReview.UserID)
			assert.Equal(t, secondReview.PieceID, createdReview.PieceID)
			assert.Equal(t, secondReview.Content, createdReview.Content)

			testReview2 = createdReview
		})

		t.Run("Test empty content no error", func(t *testing.T) {
			ctx := context.Background()

			emptyReview := testReview
			emptyReview.ID = 3
			emptyReview.Content = ""

			createdReview, err := repo.review.Create(ctx, emptyReview)
			assert.NoError(t, err)

			testReview3 = createdReview
		})
	})

	// GET ---------------------------------------------------------------------------
	t.Run("Test review get", func(t *testing.T) {
		t.Run("Test review get by ID", func(t *testing.T) {
			ctx := context.Background()

			retrievedReview, err := repo.review.GetByID(ctx, testReview.ID)
			require.NoError(t, err)
			assert.Equal(t, testReview.ID, retrievedReview.ID)
			assert.Equal(t, testReview.UserID, retrievedReview.UserID)
			assert.Equal(t, testReview.PieceID, retrievedReview.PieceID)
			assert.Equal(t, testReview.Rating, retrievedReview.Rating)
			assert.Equal(t, testReview.Content, retrievedReview.Content)
			assert.Equal(t, testReview.Moderated, retrievedReview.Moderated)
			assert.Equal(t, testReview.Published, retrievedReview.Published)
			assert.Equal(t, testReview.PhotoURL, retrievedReview.PhotoURL)
			assert.Equal(t, testReview.CreatedAt, retrievedReview.CreatedAt)
			assert.Equal(t, testReview.UpdatedAt, retrievedReview.UpdatedAt)
		})
	})

	// UPDATE ---------------------------------------------------------------------------
	t.Run("Test review update", func(t *testing.T) {
		t.Run("Test update review good", func(t *testing.T) {
			ctx := context.Background()

			updatedReview := testReview
			updatedReview.Content = "Updated review content"
			updatedReview.Rating = 8

			result, err := repo.review.Update(ctx, updatedReview)
			require.NoError(t, err)
			assert.Equal(t, updatedReview.ID, result.ID)
			assert.Equal(t, updatedReview.UserID, result.UserID)
			assert.Equal(t, updatedReview.PieceID, result.PieceID)
			assert.Equal(t, updatedReview.Rating, result.Rating)
			assert.Equal(t, updatedReview.Content, result.Content)
			assert.Equal(t, updatedReview.Moderated, result.Moderated)
			assert.Equal(t, updatedReview.Published, result.Published)
			assert.Equal(t, updatedReview.PhotoURL, result.PhotoURL)
			assert.Equal(t, updatedReview.CreatedAt, result.CreatedAt)
			assert.True(t, result.UpdatedAt.After(testReview.UpdatedAt))
		})
	})

	// DELETE ---------------------------------------------------------------------------
	t.Run("Test review delete", func(t *testing.T) {
		t.Run("Test review delete good", func(t *testing.T) {
			ctx := context.Background()

			deletedReview, err := repo.review.Delete(ctx, testReview.ID)
			require.NoError(t, err)
			assert.Equal(t, testReview.ID, deletedReview.ID)

			// Try to get the deleted review
			_, err = repo.review.GetByID(ctx, testReview.ID)
			assert.Error(t, err)

			deletedReview, err = repo.review.Delete(ctx, testReview2.ID)
			assert.NoError(t, err)
			deletedReview, err = repo.review.Delete(ctx, testReview3.ID)
			assert.NoError(t, err)

		})
	})

	// GETLIST ---------------------------------------------------------------------------
	t.Run("Test review get list", func(t *testing.T) {
		ctx := context.Background()
		// Create a few test reviews
		for i := 0; i < 5; i++ {
			reviewForSearch := domain.Review{
				ID:        i,
				UserID:    testUser.ID,
				PieceID:   uuid.New().String(),
				Rating:    5 + i,
				Content:   "Test review content " + fmt.Sprintf("%d", i),
				Moderated: i%2 == 0, // Every second review is moderated
				Published: true,
				PhotoURL:  "https://example.com/photo" + fmt.Sprintf("%d", i) + ".jpg",
			}
			_, err := repo.review.Create(ctx, reviewForSearch)
			require.NoError(t, err)
		}

		t.Run("Test review get list good", func(t *testing.T) {
			ctx := context.Background()

			pag := domain.IDPagination{
				Limit:  10,
				LastID: 0,
			}

			filter := domain.ReviewFilter{
				UserID:          &createdUser.ID,
				Published:       &[]bool{true}[0],
				IncludeProfiles: true,
			}

			reviewsReturn, newPag, err := repo.review.GetList(ctx, filter, pag)
			require.NoError(t, err)
			if reviewsReturn == nil {
				return
			}
			require.NotEmpty(t, reviewsReturn)
			assert.Equal(t, 5, len(reviewsReturn))
			assert.NotEqual(t, reviewsReturn[1].Profile.ID, uuid.Nil)
			assert.Equal(t, newPag.LastID, reviewsReturn[len(reviewsReturn)-1].ID)
		})

		t.Run("Test review get list pagination iteration", func(t *testing.T) {
			ctx := context.Background()

			pag := domain.IDPagination{
				Limit:  2,
				LastID: 0,
			}

			filter := domain.ReviewFilter{
				UserID:    &createdUser.ID,
				Published: &[]bool{true}[0],
			}

			var allReviews []domain.Review
			counter := 0
			for {
				reviews, newPag, err := repo.review.GetList(ctx, filter, pag)
				require.NoError(t, err)
				if reviews == nil || len(reviews) == 0 {
					break
				}
				if newPag.LastID == 0 {
					break // No more reviews to fetch
				}
				allReviews = append(allReviews, reviews...)
				pag = newPag
				if counter > 10 {
					break // Prevent infinite loop in case of unexpected behavior
				}
				counter++
			}

			assert.Equal(t, 5, len(allReviews))
		})

		t.Run("Test review getlist with filters", func(t *testing.T) {
			ctx := context.Background()

			pag := domain.IDPagination{
				Limit:  10,
				LastID: 0,
			}

			// Test with rating filter
			rating := 7
			filter := domain.ReviewFilter{
				UserID:    &createdUser.ID,
				Rating:    &rating,
				Published: &[]bool{true}[0],
			}

			reviews, _, err := repo.review.GetList(ctx, filter, pag)
			require.NoError(t, err)
			for _, review := range reviews {
				assert.GreaterOrEqual(t, review.Rating, rating)
			}

			// Test with moderated filter
			moderated := true
			filter = domain.ReviewFilter{
				UserID:    &createdUser.ID,
				Moderated: &moderated,
				Published: &[]bool{true}[0],
			}

			reviews, _, err = repo.review.GetList(ctx, filter, pag)
			require.NoError(t, err)
			for _, review := range reviews {
				assert.True(t, review.Moderated)
			}
		})
	})
}
