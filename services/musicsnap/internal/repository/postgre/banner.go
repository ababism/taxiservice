package postgre

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/juju/zaputil/zapctx"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"music-snap/pkg/app"
	"music-snap/pkg/msquery"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/repository/postgre/models"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
)

var _ ports.BannerRepository = &bannerRepository{}

func NewBannerRepository(db *sqlx.DB) ports.BannerRepository {
	return &bannerRepository{db: db,
		spanName: spanBaseName}
}

type bannerRepository struct {
	db       *sqlx.DB
	spanName string
}

func (r bannerRepository) ProcessDeleteQueue(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r bannerRepository) Update(initialCtx context.Context, bannerID int, banner domain.Banner) (domain.Banner, error) {
	logger := zapctx.Logger(initialCtx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(initialCtx, r.spanName+"Update")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to start transaction", err)
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	q := `
		UPDATE musicsnap b 
			SET content = $2,
			    is_active =  $3
		WHERE b.id = $1
		RETURNING *
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	writeBanner := models.ToBannerModel(banner)
	var resBanner models.Banner
	//res, err := r.db.NamedExec(q, resBanner)
	row := tx.QueryRow(q, writeBanner.Content, writeBanner.IsActive)
	if err != nil {
		_ = tx.Rollback()
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	err = row.Scan(&resBanner.ID, &resBanner.Content, &resBanner.IsActive, &resBanner.CreatedAt, &resBanner.UpdatedAt)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to create musicsnap in psql", err)

		}
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "can't scan musicsnap", err)
	}

	qFeature := `
	DELETE FROM banners_feature WHERE banner_id = $1;
	INSERT INTO banners_feature (banner_id, feature_id)
	VALUES ($1, $2);
	`

	//res, err := r.db.NamedExec(q, resBanner)
	_, err = tx.Exec(qFeature, resBanner.ID, banner.Feature)
	if err != nil {
		_ = tx.Rollback()
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	qTags := `
	DELETE FROM banners_tags WHERE banner_id = $1;
	INSERT INTO banners_tags (banner_id, tag_id)
	SELECT $1 AS banner_id, tag_id
	FROM unnest($2::int[]) AS tag_id;
	`
	tagIDs := intSliceToPostgresArray(banner.Tags)
	//res, err := r.db.NamedExec(q, resBanner)
	_, err = tx.Exec(qTags, resBanner.ID, tagIDs)
	if err != nil {
		_ = tx.Rollback()
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to commit transaction", err)
	}

	return resBanner.ToDomain(), nil
}

func (r bannerRepository) Find(initialCtx context.Context, tag, feature int) (domain.Banner, error) {
	logger := zapctx.Logger(initialCtx)

	tr := global.Tracer(domain.ServiceName)
	_, span := tr.Start(initialCtx, r.spanName+"Find")
	defer span.End()

	q := `
	SELECT b.* FROM musicsnap b
	    JOIN banners_feature bf on b.id = bf.banner_id 
	    JOIN public.banners_tags bt on b.id = bt.banner_id
	WHERE  bf.feature_id = $1 AND bt.tag_id = $2
	LIMIT 1
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var banner models.Banner
	err := r.db.Get(&banner, q, feature, tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Banner{}, app.NewError(http.StatusNotFound, "musicsnap not found", "failed to get musicsnap from DB", err)
		}
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}
	return banner.ToDomain(), nil
}

// TODO Add features and tags
func (r bannerRepository) Get(initialCtx context.Context, ID int) (domain.Banner, error) {
	logger := zapctx.Logger(initialCtx)

	tr := global.Tracer(domain.ServiceName)
	_, span := tr.Start(initialCtx, r.spanName+"UpdatePr")
	defer span.End()

	q := `
	SELECT * FROM musicsnap WHERE id = $1
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var banner models.Banner
	err := r.db.Get(&banner, q, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Banner{}, app.NewError(http.StatusNotFound, "musicsnap not found", "failed to get musicsnap from DB", err)
		}
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}
	return banner.ToDomain(), nil
}

func (r bannerRepository) Create(initialCtx context.Context, banner domain.Banner) (domain.Banner, error) {
	logger := zapctx.Logger(initialCtx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(initialCtx, r.spanName+"UpdatePr")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to start transaction", err)
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	q := `
	INSERT INTO musicsnap (content, is_active)
	VALUES ($1, $2)
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	writeBanner := models.ToBannerModel(banner)
	var resBanner models.Banner
	//res, err := r.db.NamedExec(q, resBanner)
	row := tx.QueryRow(q, writeBanner.Content, writeBanner.IsActive)
	if err != nil {
		_ = tx.Rollback()
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	err = row.Scan(&resBanner.ID, &resBanner.Content, &resBanner.IsActive, &resBanner.CreatedAt, &resBanner.UpdatedAt)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to create musicsnap in psql", err)

		}
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "can't scan musicsnap", err)
	}

	qFeature := `
	INSERT INTO banners_feature (banner_id, feature_id)
	VALUES ($1, $2);
	`

	//res, err := r.db.NamedExec(q, resBanner)
	_, err = tx.Exec(qFeature, resBanner.ID, banner.Feature)
	if err != nil {
		_ = tx.Rollback()
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	qTags := `
	INSERT INTO banners_tags (banner_id, tag_id)
	SELECT $1 AS banner_id, tag_id
	FROM unnest($2::int[]) AS tag_id;
	`
	tagIDs := intSliceToPostgresArray(banner.Tags)
	//res, err := r.db.NamedExec(q, resBanner)
	_, err = tx.Exec(qTags, resBanner.ID, tagIDs)
	if err != nil {
		_ = tx.Rollback()
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	err = tx.Commit()
	if err != nil {
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to commit transaction", err)
	}

	return resBanner.ToDomain(), nil
}

func mapBannerFilterParams(params *domain.BannerFilter) map[string]any {
	if params == nil {
		return map[string]any{}
	}
	paramsMap := make(map[string]any)
	if params.TagID != nil {
		paramsMap["tag_id"] = params.TagID
	}
	if params.Feature != nil {
		paramsMap["feature_id"] = params.Feature
	}
	return paramsMap
}

// TODO Refactor?
func (r bannerRepository) GetList(initialCtx context.Context, filter domain.BannerFilter) ([]domain.Banner, error) {
	logger := zapctx.Logger(initialCtx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(initialCtx, r.spanName+"UpdatePr")
	defer span.End()

	baseQuery := `
	SELECT b.*, bf.feature_id, bt.tag_id FROM musicsnap b
	    JOIN banners_feature bf on b.id = bf.banner_id
	    JOIN public.banners_tags bt on b.id = bt.banner_id
	`

	paramsMap := mapBannerFilterParams(&filter)

	qb := msquery.NewFilterBuilder(baseQuery, 0)
	query, args, _ := qb.
		WhereAnd(paramsMap).
		Offset(*filter.Offset).
		Limit(*filter.Limit).
		OrderBy("created_at", true).
		Build()

	logger.With(zap.String("PSQL query", formatQuery(query)))

	type Row struct {
		models.Banner
		FeatureID int `db:"feature_id"`
		TagID     int `db:"tag_id"`
	}
	var banners []Row
	err := r.db.SelectContext(ctx, &banners, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.NewError(http.StatusNotFound, "musicsnap not found", "failed to get musicsnap from DB", err)
		}
		return nil, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)

	}

	domainBanners := make([]domain.Banner, 0, len(banners))
	for _, banner := range banners {
		// TODO get all tags
		banner.Tags = []int{banner.TagID}
		banner.Feature = banner.FeatureID
		domainBanners = append(domainBanners, banner.ToDomain())
	}

	return domainBanners, nil
}

// TODO Remove?
//func (r bannerRepository) GetListNew(initialCtx context.Context, filter domain.BannerFilter) ([]domain.Banner, error) {
//	logger := zapctx.Logger(initialCtx)
//
//	tr := global.Tracer(domain.ServiceName)
//	ctx, span := tr.Start(initialCtx, r.spanName+"Create")
//	defer span.End()
//
//	baseQuery := `
//	SELECT b.id, content, is_active, created_at, updated_at, bf.feature_id, bt.tag_id FROM musicsnap b
//	    JOIN banners_feature bf on b.id = bf.banner_id
//	    JOIN public.banners_tags bt on b.id = bt.banner_id
//	`
//
//	paramsMap := mapBannerFilterParams(&filter)
//
//	qb := msquery.NewFilterBuilder(baseQuery, 0)
//	query, args, _ := qb.
//		WhereAnd(paramsMap).
//		Offset(*filter.Offset).
//		Limit(*filter.Limit).
//		OrderBy("created_at", true).
//		Build()
//
//	logger.With(zap.String("PSQL query", formatQuery(query)))
//
//	type Row struct {
//		models.Banner
//		FeatureID int `db:"feature_id"`
//		TagID     int `db:"tag_id"`
//	}
//
//	rows, err := r.db.Query(query, args...)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, app.NewError(http.StatusNotFound, "musicsnap not found", "failed to get musicsnap from DB", err)
//		}
//		return nil, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//
//	}
//	var row Row
//	var banner models.Banner
//	var pbID int
//	domainBanners := make([]domain.Banner, 0)
//	for {
//		rows.Next()
//		err := rows.Scan(&row)
//		if err != nil {
//			if errors.Is(err, sql.ErrNoRows) {
//				break
//			}
//			return nil, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//		}
//		// TODO FINISH
//		banner = models.Banner{
//			ID:        row.Banner.ID,
//			Content:   row.Banner.Content,
//			IsActive:  row.Banner.IsActive,
//			Feature:   row.FeatureID,
//			Tags:      nil,
//			CreatedAt: row.Banner.CreatedAt,
//			UpdatedAt: row.Banner.UpdatedAt,
//		}
//		if row.Banner.ID != pbID {
//			banner.Tags = []int{row.TagID}
//		} else {
//			banner.Tags = append(banner.Tags, row.TagID)
//		}
//		banner.Feature = row.FeatureID
//		domainBanners = append(domainBanners, banner.ToValidDomain())
//
//	}
//	return domainBanners, nil
//}
//
//func (r bannerRepository) ProcessDeleteQueue(ctx context.Context) error {
//	//TODO implement me
//	panic("implement me")
//}

func (r bannerRepository) UpdateWithNoRelations(ctx context.Context, bannerID int, banner domain.Banner) (domain.Banner, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"UpdatePr")
	defer span.End()

	// TODO ADD TX to update relations
	q := `
		UPDATE musicsnap b 
			SET content = $2,
			    is_active =  $3
		WHERE b.id = $1
		RETURNING *
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var res models.Banner
	err := r.db.Get(&res, q, banner.ID, banner.Content, banner.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Banner{}, app.NewError(http.StatusNotFound, "musicsnap not found", "failed to get musicsnap from DB", err)
		}
		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}
	return res.ToDomain(), nil
}

func (r bannerRepository) Delete(ctx context.Context, bannerID int) error {
	//TODO implement me
	panic("implement me")
}

//func (r bannerRepository) Update(initialCtx context.Context, banner domain.Banner) (domain.Banner, error) {
//	logger := zapctx.Logger(initialCtx)
//
//	tr := global.Tracer(domain.ServiceName)
//	ctx, span := tr.Start(initialCtx, r.spanName+"Update")
//	defer span.End()
//
//	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
//	if err != nil {
//		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to start transaction", err)
//	}
//
//	defer func(tx *sql.Tx) {
//		_ = tx.Rollback()
//	}(tx)
//
//	q := `
//		UPDATE musicsnap b
//			SET content = $2,
//			    is_active =  $3
//		WHERE b.id = $1
//		RETURNING *
//	`
//	logger.With(zap.String("PSQL query", formatQuery(q)))
//
//	writeBanner := models.ToBannerModel(banner)
//	var resBanner models.Banner
//	//res, err := r.db.NamedExec(q, resBanner)
//	row := tx.QueryRow(q, writeBanner.Content, writeBanner.IsActive)
//	if err != nil {
//		_ = tx.Rollback()
//		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//	}
//
//	err = row.Scan(&resBanner.ID, &resBanner.Content, &resBanner.IsActive, &resBanner.CreatedAt, &resBanner.UpdatedAt)
//	if err != nil {
//		_ = tx.Rollback()
//		if errors.Is(err, sql.ErrNoRows) {
//			return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to create musicsnap in psql", err)
//
//		}
//		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "can't scan musicsnap", err)
//	}
//
//	qFeature := `
//	DELETE FROM banners_feature WHERE banner_id = $1;
//	INSERT INTO banners_feature (banner_id, feature_id)
//	VALUES ($1, $2);
//	`
//
//	//res, err := r.db.NamedExec(q, resBanner)
//	_, err = tx.Exec(qFeature, resBanner.ID, banner.Feature)
//	if err != nil {
//		_ = tx.Rollback()
//		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//	}
//
//	qTags := `
//	DELETE FROM banners_tags WHERE banner_id = $1;
//	INSERT INTO banners_tags (banner_id, tag_id)
//	SELECT $1 AS banner_id, tag_id
//	FROM unnest($2::int[]) AS tag_id;
//	`
//	tagIDs := intSliceToPostgresArray(banner.Tags)
//	//res, err := r.db.NamedExec(q, resBanner)
//	_, err = tx.Exec(qTags, resBanner.ID, tagIDs)
//	if err != nil {
//		_ = tx.Rollback()
//		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//	}
//
//	err = tx.Commit()
//	if err != nil {
//		return domain.Banner{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to commit transaction", err)
//	}
//
//	return resBanner.ToValidDomain(), nil
//}
