package postgre

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"music-snap/services/musicsnap/internal/service/ports"
	"strings"
)

const (
	spanBaseName = "musicsnap/repository/postgres."
)

type Repository struct {
	User     ports.UserRepository
	Review   ports.ReviewRepository
	Reaction ports.ReactionRepository
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{
		User:     NewUserRepository(db),
		Review:   NewReviewRepository(db),
		Reaction: NewReactionRepository(db),
	}
}

type repository struct {
	user     userRepository
	review   reviewRepository
	reaction reactionRepository
}

func newRepository(db *sqlx.DB) repository {
	return repository{
		user:     newUserRepository(db),
		review:   newReviewRepository(db),
		reaction: newReactionRepository(db),
	}
}

const (
	UserTable   = "users"
	ReviewTable = "reviews"
)

func formatQuery(q string) string {
	return fmt.Sprintf("SQL Query: %s", strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " "))
}

func intSliceToPostgresArray(slice []int) string {
	return "{" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(slice)), ","), "[]") + "}"
}
func stringSliceToPostgresArray(slice []string) string {
	return "{" +
		strings.Trim(
			strings.Join(
				strings.Fields(fmt.Sprint(slice)), ","), "[]") +
		"}"
}
