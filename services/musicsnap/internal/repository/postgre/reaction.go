package postgre

import (
	c "context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"music-snap/pkg/app"

	//"database/sql"
	//"github.com/google/uuid"
	//"github.com/jmoiron/sqlx"
	"github.com/juju/zaputil/zapctx"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	//"music-snap/pkg/app"
	//"music-snap/pkg/msquery"
	//"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/repository/postgre/models"
	//"music-snap/services/musicsnap/internal/service/ports"
	//"net/http"
	"github.com/jmoiron/sqlx"
	domain "music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
)

var _ ports.ReactionRepository = &reactionRepository{}

func NewReactionRepository(db *sqlx.DB) ports.ReactionRepository {
	return &reactionRepository{db: db,
		spanName: spanBaseName + "reactionRepository."}
}

func newReactionRepository(db *sqlx.DB) reactionRepository {
	return reactionRepository{db: db,
		spanName: spanBaseName + "reactionRepository."}
}

type reactionRepository struct {
	db       *sqlx.DB
	spanName string
}

func (r reactionRepository) GetFromActor(ctx c.Context, userID uuid.UUID, reviewID int) (domain.Reaction, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetFromActor")
	defer span.End()

	q := `
	SELECT * FROM reactions
	WHERE user_id = $1 AND review_id = $2;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var reaction models.ReactionModel
	err := r.db.GetContext(ctx, &reaction, q, userID, reviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Reaction{}, app.NewError(http.StatusNotFound, "reaction not found", "reaction with given user and review does not exist", err)
		}
		return domain.Reaction{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return reaction.ToDomain(), nil
}

func (r reactionRepository) Create(ctx c.Context, reaction domain.Reaction) (domain.Reaction, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"Create")
	defer span.End()

	q := `
	INSERT INTO reactions (user_id, review_id, type)
	VALUES ($1, $2, $3)
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	reactionToWrite := models.ToReactionModel(reaction)

	var createdReaction models.ReactionModel
	err := r.db.GetContext(ctx, &createdReaction, q, reactionToWrite.UserID, reactionToWrite.ReviewID, reactionToWrite.Type)
	if err != nil {
		return domain.Reaction{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return createdReaction.ToDomain(), nil
}

func (r reactionRepository) GetByID(ctx c.Context, id int) (domain.Reaction, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetByID")
	defer span.End()

	q := `
	SELECT * FROM reactions
	WHERE id = $1;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var reaction models.ReactionModel
	err := r.db.GetContext(ctx, &reaction, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Reaction{}, app.NewError(http.StatusNotFound, "reaction not found", "reaction with given id does not exist", err)
		}
		return domain.Reaction{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return reaction.ToDomain(), nil
}

func (r reactionRepository) GetFromReview(ctx c.Context, reviewID int) ([]domain.Reaction, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetFromReview")
	defer span.End()

	q := `
	SELECT * FROM reactions
	WHERE review_id = $1;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var reactions []models.ReactionModel
	err := r.db.SelectContext(ctx, &reactions, q, reviewID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Reaction{}, nil
		}
		return nil, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	reactionsRes := make([]domain.Reaction, 0, len(reactions))
	for _, r := range reactions {
		reactionsRes = append(reactionsRes, r.ToDomain())
	}

	return reactionsRes, nil
}

func (r reactionRepository) Update(ctx c.Context, reaction domain.Reaction) (domain.Reaction, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"Update")
	defer span.End()

	q := `
	UPDATE reactions
	SET type = $1
	WHERE id = $2
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	reactionToWrite := models.ToReactionModel(reaction)

	var updatedReaction models.ReactionModel
	err := r.db.GetContext(ctx, &updatedReaction, q, reactionToWrite.Type, reactionToWrite.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Reaction{}, app.NewError(http.StatusNotFound, "reaction not found", "reaction with given id does not exist", err)
		}
		return domain.Reaction{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return updatedReaction.ToDomain(), nil
}

func (r reactionRepository) Delete(ctx c.Context, id int) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"Delete")
	defer span.End()

	q := `
	DELETE FROM reactions
	WHERE id = $1
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var reaction models.ReactionModel
	err := r.db.GetContext(ctx, &reaction, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.NewError(http.StatusNotFound, "reaction not found", "reaction with given id does not exist", err)
		}
		return app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return nil
}
