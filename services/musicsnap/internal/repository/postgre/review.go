package postgre

import (
	c "context"
	"database/sql"
	"errors"
	"github.com/juju/zaputil/zapctx"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"music-snap/pkg/app"
	qb "music-snap/pkg/querybuilder"
	"music-snap/services/musicsnap/internal/repository/postgre/models"
	"net/http"

	//"database/sql"
	//"github.com/google/uuid"
	//"github.com/jmoiron/sqlx"
	//"github.com/juju/zaputil/zapctx"
	//global "go.opentelemetry.io/otel"
	//"go.uber.org/zap"
	//"music-snap/pkg/app"
	//"music-snap/pkg/msquery"
	//"music-snap/services/musicsnap/internal/domain"
	//"music-snap/services/musicsnap/internal/repository/postgre/models"
	//"music-snap/services/musicsnap/internal/service/ports"
	//"net/http"
	"github.com/jmoiron/sqlx"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
)

var _ ports.ReviewRepository = &reviewRepository{}

func NewReviewRepository(db *sqlx.DB) ports.ReviewRepository {
	return &reviewRepository{db: db,
		spanName: spanBaseName + "reviewRepository."}
}

func newReviewRepository(db *sqlx.DB) reviewRepository {
	return reviewRepository{db: db,
		spanName: spanBaseName + "reviewRepository."}
}

type reviewRepository struct {
	db       *sqlx.DB
	spanName string
}

func (r reviewRepository) Create(ctx c.Context, review domain.Review) (domain.Review, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"Create")
	defer span.End()

	q := `
	INSERT INTO reviews (user_id, piece_id, rating, photo_url, content, moderated, published)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	reviewToWrite := models.ToReviewModel(review)

	var createdReview models.ReviewModel
	err := r.db.GetContext(ctx, &createdReview, q, reviewToWrite.UserID, reviewToWrite.PieceID, reviewToWrite.Rating, reviewToWrite.PhotoURL, reviewToWrite.Content, reviewToWrite.Moderated, reviewToWrite.Published)
	if err != nil {
		return domain.Review{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return createdReview.ToDomain(), nil
}

func (r reviewRepository) Update(ctx c.Context, review domain.Review) (domain.Review, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"Update")
	defer span.End()

	q := `
	UPDATE reviews
	SET piece_id = $1, rating = $2, photo_url = $3, content = $4, moderated = $5, published = $6
	WHERE id = $7
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	reviewToWrite := models.ToReviewModel(review)

	var updatedReview models.ReviewModel
	err := r.db.GetContext(ctx, &updatedReview, q, reviewToWrite.PieceID, reviewToWrite.Rating, reviewToWrite.PhotoURL, reviewToWrite.Content, reviewToWrite.Moderated, reviewToWrite.Published, review.ID)
	if err != nil {
		return domain.Review{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return updatedReview.ToDomain(), nil
}

func (r reviewRepository) GetByID(ctx c.Context, id int) (domain.Review, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetByID")
	defer span.End()

	q := `
	SELECT * FROM reviews
	WHERE id = $1;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var review models.ReviewModel
	err := r.db.GetContext(ctx, &review, q, id)
	if err != nil {
		return domain.Review{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return review.ToDomain(), nil
}

func (r reviewRepository) GetList(ctx c.Context, filter domain.ReviewFilter, pag domain.IDPagination) ([]domain.Review, domain.IDPagination, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetList")
	defer span.End()

	//orderByField := "created_at"
	//if filter.OrderByRating != nil && *filter.OrderByRating {
	//	orderByField = "rating"
	//}
	orderByField := "reviews.id"
	filter.OrderAsc = true

	// TODO Add OPT join
	qBuild := qb.NewNamed().
		Q("SELECT reviews.*, users.id AS user_id, users.nickname,users.avatar_url, users.background_url, users.bio  FROM reviews").
		StartOpt().
		Q("JOIN").Table("users").ON().Q("reviews.user_id = users.id").
		EndOptIf(func() bool {
			return true
		}).
		WhereOptPart().
		CompConnectorOpt("user_id", qb.EQ(), "user_id", filter.UserID, qb.AND()).
		CompConnectorOpt("piece_id", qb.EQ(), "piece_id", filter.PieceID, qb.AND()).
		CompConnectorOpt("rating", qb.GET(), "rating", filter.Rating, qb.AND()).
		CompConnectorOpt("moderated", qb.EQ(), "moderated", filter.Moderated, qb.AND()).
		CompConnectorOpt("published", qb.EQ(), "published", filter.Published, qb.AND()).
		CompConnectorOpt("reviews.id", qb.GT(), "last_id", pag.LastID, qb.AND()).
		EndWhereOpt().
		OrderBy(orderByField, filter.OrderAsc).
		Limit("", 10)
	q, args := qBuild.Build()

	logger.With(zap.String("PSQL query", formatQuery(q)))

	logger.Warn("Executing query", zap.String("query", q), zap.Any("args", args))

	type ReviewAndProfileRow struct {
		models.ReviewModel
		models.UserModel
	}

	var reviewRows []ReviewAndProfileRow
	preparedQ, err := r.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, pag, app.NewError(http.StatusInternalServerError, "unknown error", "internal error preparing named query", err)
	}
	err = preparedQ.SelectContext(ctx, &reviewRows, args)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pag, nil // No reviews found
		}
		return nil, pag, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	reviewsRes := make([]domain.Review, 0, len(reviewRows))
	if filter.IncludeProfiles {
		for _, rRow := range reviewRows {
			reviewsRes = append(reviewsRes, rRow.ReviewModel.ToDomainProfile(rRow.UserModel.ToProfileDomain()))
		}
	} else {
		for _, rRow := range reviewRows {
			reviewsRes = append(reviewsRes, rRow.ReviewModel.ToDomain())
		}
	}

	if len(reviewsRes) == 0 {
		pag.LastID = 0
		return []domain.Review{}, pag, nil
	}

	newLastID := reviewsRes[len(reviewsRes)-1].ID
	pag.LastID = newLastID

	return reviewsRes, pag, nil
}

func (r reviewRepository) Delete(ctx c.Context, id int) (domain.Review, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"Delete")
	defer span.End()

	q := `
	DELETE FROM reviews
	WHERE id = $1
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var review models.ReviewModel
	err := r.db.GetContext(ctx, &review, q, id)
	if err != nil {
		return domain.Review{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return review.ToDomain(), nil
}
