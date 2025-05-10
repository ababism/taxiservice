package postgre

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
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

var _ ports.UserRepository = &userRepository{}

func NewUserRepository(db *sqlx.DB) ports.UserRepository {
	return &userRepository{db: db,
		spanName: spanBaseName + "userRepository."}
}

func newUserRepository(db *sqlx.DB) userRepository {
	return userRepository{db: db,
		spanName: spanBaseName + "userRepository."}
}

type userRepository struct {
	db       *sqlx.DB
	spanName string
}

func (r userRepository) GetList(ctx context.Context, nickQuery string, pag domain.UUIDPagination) ([]domain.User, uuid.UUID, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetList")
	defer span.End()

	q := `
	SELECT * FROM users 
	         WHERE nickname ILIKE $1 AND id > $2
	ORDER BY id ASC 
	LIMIT $3;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var resUsers []models.UserModel
	err := r.db.SelectContext(ctx, &resUsers, q, "%"+nickQuery+"%", pag.LastUUID, pag.Limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.User{}, uuid.Nil, app.NewError(http.StatusNotFound, "user not found", "user not found", err)
		}
		return []domain.User{}, uuid.Nil, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	var users []domain.User
	for _, user := range resUsers {
		q = `
		SELECT * FROM user_roles WHERE user_id = $1 ORDER BY role ASC;
		`
		var roles []models.RoleModel
		err = r.db.SelectContext(ctx, &roles, q, user.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return []domain.User{}, uuid.Nil, app.NewError(http.StatusNotFound, "user roles not found", "user roles not found", err)
			}
			return []domain.User{}, uuid.Nil, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
		}
		users = append(users, user.ToDomain(roles))
	}

	if len(users) == 0 {
		return []domain.User{}, uuid.Nil, nil
	}

	newLastUUID := users[len(users)-1].ID
	if len(users) == 0 {
		newLastUUID = uuid.Nil
	}
	return users, newLastUUID, nil
}

func (r userRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"UpdatePr")
	defer span.End()

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to start transaction", err)
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	q := `
	INSERT INTO users (id, nickname, avatar_url, background_url, bio, email, password_hash)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	writeUser := models.ToUserModel(user)

	if writeUser.ID == uuid.Nil {
		writeUser.ID = uuid.New()
	}

	var resUser models.UserModel
	//res, err := r.db.NamedExec(q, resUser)
	row := tx.QueryRow(q, writeUser.ID, writeUser.Nickname, writeUser.AvatarURL, writeUser.BackgroundURL, writeUser.Bio, writeUser.Email, writeUser.PasswordHash)
	if err != nil {
		_ = tx.Rollback()
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	err = row.Scan(&resUser.ID, &resUser.Nickname, &resUser.AvatarURL, &resUser.BackgroundURL, &resUser.Bio, &resUser.Email, &resUser.PasswordHash, &resUser.CreatedAt, &resUser.UpdatedAt)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to create user in psql", err)

		}
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "can't scan user", err)
	}

	//qRolesRange := `
	//INSERT INTO user_roles (user_id, role)
	//VALUES ($1, $2);
	//`
	//for role := range user.Roles.ToSlice() {
	//	_, err = tx.Exec(qRolesRange, writeUser.ID, role)
	//	if err != nil {
	//		_ = tx.Rollback()
	//		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	//	}
	//}
	//res, err := r.db.NamedExec(q, resUser)

	qRoles := `
	INSERT INTO user_roles (user_id, role)
	SELECT $1 AS user_id, role
	FROM unnest($2::text[]) AS role;
	`
	rolesIDs := stringSliceToPostgresArray(user.Roles.ToSlice())
	//res, err := r.db.NamedExec(q, resUser)
	_, err = tx.Exec(qRoles, resUser.ID, rolesIDs)
	if err != nil {
		_ = tx.Rollback()
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	err = tx.Commit()
	if err != nil {
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "failed to commit transaction", err)
	}

	q = `
	SELECT * FROM user_roles WHERE user_id = $1 ORDER BY role ASC;
	`
	var roles []models.RoleModel
	err = r.db.SelectContext(ctx, &roles, q, resUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, app.NewError(http.StatusNotFound, "user roles not found", "user roles not found", err)
		}
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return resUser.ToDomain(roles), nil
}

func mapUserFilterParams(params *domain.UserFilter) map[string]any {
	if params == nil {
		return map[string]any{}
	}
	paramsMap := make(map[string]any)
	if params.ID != nil && *params.ID != uuid.Nil {
		paramsMap["id"] = params.ID
	}
	if params.Nickname != nil {
		paramsMap["nickname"] = params.Nickname
	}
	if params.Email != nil {
		paramsMap["email"] = params.Email
	}
	return paramsMap
}

func (r userRepository) GetByParam(ctx context.Context, filter domain.UserFilter) (domain.User, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetByID")
	defer span.End()

	baseQuery := `
	SELECT * FROM users
	`
	paramsMap := mapUserFilterParams(&filter)

	qb := msquery.NewFilterBuilder(baseQuery, 0)
	query, args, _ := qb.
		WhereAnd(paramsMap).
		Limit(1).
		Build()

	logger.With(zap.String("PSQL query", formatQuery(query)))
	//fmt.Println("PSQL QUERY PARAM" + formatQuery(query))

	var resUser models.UserModel
	err := r.db.GetContext(ctx, &resUser, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, app.NewError(http.StatusNotFound, "user not found", "user not found", err)
		}
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}
	q := `
	SELECT * FROM user_roles WHERE user_id = $1 ORDER BY created_at ASC;
	`
	var roles []models.RoleModel
	err = r.db.SelectContext(ctx, &roles, q, resUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, app.NewError(http.StatusNotFound, "user roles not found", "user roles not found", err)
		}
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return resUser.ToDomain(roles), nil
}

func (r userRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetByID")
	defer span.End()

	q := `
	SELECT * FROM users WHERE id = $1;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var resUser models.UserModel
	err := r.db.GetContext(ctx, &resUser, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, app.NewError(http.StatusNotFound, "user not found", "user not found", err)
		}
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}
	q = `
	SELECT * FROM user_roles WHERE user_id = $1 ORDER BY role ASC;
	`
	var roles []models.RoleModel
	err = r.db.SelectContext(ctx, &roles, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, app.NewError(http.StatusNotFound, "user roles not found", "user roles not found", err)
		}
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return resUser.ToDomain(roles), nil
}

func (r userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.GetByParam(ctx, domain.UserFilter{Email: &email})
}

//func (r userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
//	logger := zapctx.Logger(ctx)
//
//	tr := global.Tracer(domain.ServiceName)
//	ctx, span := tr.Start(ctx, r.spanName+"GetByEmail")
//	defer span.End()
//
//	q := `
//	SELECT * FROM users WHERE email = $1;
//	`
//	logger.With(zap.String("PSQL query", formatQuery(q)))
//
//	var resUser models.UserModel
//	err := r.db.GetContext(ctx, &resUser, q, email)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return domain.User{}, app.NewError(http.StatusNotFound, "user not found", "user not found", err)
//		}
//		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//	}
//
//	q = `
//	SELECT * FROM user_roles WHERE user_id = $1 ORDER BY role ASC;
//	`
//	var roles []models.RoleModel
//	err = r.db.SelectContext(ctx, &roles, q, resUser.ID)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return domain.User{}, app.NewError(http.StatusNotFound, "user roles not found", "user roles not found", err)
//		}
//		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//	}
//
//	return resUser.ToValidDomain(roles), nil
//}

func (r userRepository) GetByNickname(ctx context.Context, nickname string) (domain.User, error) {
	return r.GetByParam(ctx, domain.UserFilter{Nickname: &nickname})
}

//func (r userRepository) GetByNickname(ctx context.Context, nickname string) (domain.User, error) {
//	logger := zapctx.Logger(ctx)
//
//	tr := global.Tracer(domain.ServiceName)
//	ctx, span := tr.Start(ctx, r.spanName+"GetByNickname")
//	defer span.End()
//
//	q := `
//	SELECT * FROM users WHERE nickname = $1;
//	`
//	logger.With(zap.String("PSQL query", formatQuery(q)))
//
//	var resUser models.UserModel
//	err := r.db.GetContext(ctx, &resUser, q, nickname)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return domain.User{}, app.NewError(http.StatusNotFound, "user not found", "user not found", err)
//		}
//		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//	}
//	q = `
//	SELECT * FROM user_roles WHERE user_id = $1 ORDER BY role ASC;
//	`
//	var roles []models.RoleModel
//	err = r.db.SelectContext(ctx, &roles, q, resUser.ID)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return domain.User{}, app.NewError(http.StatusNotFound, "user roles not found", "user roles not found", err)
//		}
//		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
//	}
//
//	return resUser.ToValidDomain(roles), nil
//}

func (r userRepository) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"UpdateUser")
	defer span.End()

	q := `
	UPDATE users SET (nickname, avatar_url, background_url, bio, email, password_hash) = ($1, $2, $3, $4, $5, $6)
	WHERE id = $7
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	writeUser := models.ToUserModel(user)
	resUser := models.UserModel{}
	err := r.db.GetContext(ctx, &resUser, q, writeUser.Nickname, writeUser.AvatarURL, writeUser.BackgroundURL, writeUser.Bio, writeUser.Email, writeUser.PasswordHash, writeUser.ID)
	if err != nil {
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}
	//r.db.GetContext(ctx, &resUser, q, writeUser.ID)
	q = `
	SELECT * FROM user_roles WHERE user_id = $1 ORDER BY role;
	`
	var roles []models.RoleModel
	err = r.db.SelectContext(ctx, &roles, q, resUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, app.NewError(http.StatusNotFound, "user roles not found", "user roles not found", err)
		}
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return resUser.ToDomain(roles), nil
}

func (r userRepository) UpdateProfile(ctx context.Context, user domain.User) (domain.User, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"UpdateUser")
	defer span.End()

	q := `
	UPDATE users SET (nickname, avatar_url, background_url, bio) = ($1, $2, $3, $4)
	WHERE id = $5
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	writeUser := models.ToUserModel(user)
	resUser := models.UserModel{}
	err := r.db.GetContext(ctx, &resUser, q, writeUser.Nickname, writeUser.AvatarURL, writeUser.BackgroundURL, writeUser.Bio, writeUser.ID)
	if err != nil {
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}
	//r.db.GetContext(ctx, &resUser, q, writeUser.ID)
	q = `
	SELECT * FROM user_roles WHERE user_id = $1 ORDER BY role;
	`
	var roles []models.RoleModel
	err = r.db.SelectContext(ctx, &roles, q, resUser.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, app.NewError(http.StatusNotFound, "user roles not found", "user roles not found", err)
		}
		return domain.User{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return resUser.ToDomain(roles), nil
}

func (r userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"Delete")
	defer span.End()
	q := `
	DELETE FROM user_roles WHERE user_id = $1;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	q = `
	DELETE FROM users WHERE id = $1;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	_, err = r.db.ExecContext(ctx, q, id)
	if err != nil {
		return app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return nil
}

func (r userRepository) CreateSub(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"CreateSub")
	defer span.End()

	q := `
	INSERT INTO subscriptions (subscriber_id, followed_id, notification_flag)
	VALUES ($1, $2, $3)
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	writeSub := models.ToSubscriptionModel(sub)

	var resSub models.SubscriptionModel
	err := r.db.GetContext(ctx, &resSub, q, writeSub.SubscriberID, writeSub.FollowedID, writeSub.NotificationFlag)
	if err != nil {
		return domain.Subscription{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return resSub.ToLightDomain(), nil
}

func (r userRepository) GetSub(ctx context.Context, subscriberID uuid.UUID, followedID uuid.UUID) (domain.Subscription, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetSub")
	defer span.End()

	q := `
	SELECT * FROM subscriptions WHERE subscriber_id = $1 AND followed_id = $2;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var resSub models.SubscriptionModel
	err := r.db.GetContext(ctx, &resSub, q, subscriberID, followedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Subscription{}, app.NewError(http.StatusNotFound, "subscription not found", "subscription not found", err)
		}
		return domain.Subscription{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return resSub.ToLightDomain(), nil
}

func (r userRepository) UpdateSub(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"UpdateSub")
	defer span.End()

	q := `
	UPDATE subscriptions SET (notification_flag) = ($1)
	WHERE id = $2
	RETURNING *;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	writeSub := models.ToSubscriptionModel(sub)

	var resSub models.SubscriptionModel
	err := r.db.GetContext(ctx, &resSub, q, writeSub.NotificationFlag, writeSub.ID)
	if err != nil {
		return domain.Subscription{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return resSub.ToLightDomain(), nil
}

func (r userRepository) DeleteSub(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"DeleteSub")
	defer span.End()

	q := `
	DELETE FROM subscriptions WHERE sub_id = $1;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	_, err := r.db.ExecContext(ctx, q, sub.ID)
	if err != nil {
		return domain.Subscription{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	return domain.Subscription{}, nil
}

func (r userRepository) ListSubscriptions(ctx context.Context, subscriberID uuid.UUID, followedID uuid.UUID, pag domain.IDPagination) ([]domain.Subscription, domain.IDPagination, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"ListSubscriptions")
	defer span.End()

	var q string
	if subscriberID != uuid.Nil {
		q = `
			SELECT *
			FROM subscriptions s
					JOIN users u 
					    ON s.followed_id = u.id 
					           AND s.subscriber_id = $1
			WHERE s.sub_id > $2
				ORDER BY s.sub_id ASC
			LIMIT $3;
			`
	} else if followedID != uuid.Nil {
		q = `
			SELECT s.*, u.*
			FROM subscriptions s
					JOIN users u 
					    ON s.subscriber_id = u.id 
					           AND s.followed_id = $1
			WHERE s.sub_id > $2
				ORDER BY s.sub_id ASC
			LIMIT $3;
			`
	} else {
		return []domain.Subscription{}, pag, app.NewError(http.StatusBadRequest, "no targeted user ID in request", "no subscriberID or followedID in method", nil)
	}

	logger.With(zap.String("PSQL query", formatQuery(q)))

	type Row struct {
		models.SubscriptionModel
		models.UserModel
	}

	var resSubs []Row
	var err error

	if followedID != uuid.Nil {
		err = r.db.SelectContext(ctx, &resSubs, q, followedID, pag.LastID, pag.Limit)
	} else {
		err = r.db.SelectContext(ctx, &resSubs, q, subscriberID, pag.LastID, pag.Limit)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Subscription{}, pag, app.NewError(http.StatusNotFound, "subscriptions not found", "subscriptions not found", err)
		}
		return []domain.Subscription{}, pag, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	subs := make([]domain.Subscription, 0, len(resSubs))
	for _, sub := range resSubs {
		subs = append(subs, sub.SubscriptionModel.ToDomain(sub.UserModel.ToProfileDomain()))
	}

	if len(subs) == 0 {
		pag.LastID = 0
		return []domain.Subscription{}, pag, nil
	}

	newLastID := subs[len(subs)-1].ID
	pag.LastID = newLastID

	return subs, pag, nil
}

func (r userRepository) GetProfile(ctx context.Context, profileID uuid.UUID) (domain.Profile, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetProfile")
	defer span.End()

	q := `
	SELECT * FROM users WHERE id = $1;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var resUser models.UserModel
	err := r.db.GetContext(ctx, &resUser, q, profileID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Profile{}, app.NewError(http.StatusNotFound, "Profile not found", "user not found", err)
		}
		return domain.Profile{}, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}
	return resUser.ToProfileDomain(), nil
}

func (r userRepository) GetProfileListByUUIDSlice(ctx context.Context, profileIDs []uuid.UUID, pag domain.UUIDPagination) ([]domain.Profile, uuid.UUID, error) {
	logger := zapctx.Logger(ctx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, r.spanName+"GetList")
	defer span.End()

	q := `
	SELECT * FROM users 
	         WHERE id IN $1 AND id > $2
	ORDER BY id ASC 
	LIMIT $3;
	`
	logger.With(zap.String("PSQL query", formatQuery(q)))

	var extractedProfiles []models.UserModel
	err := r.db.SelectContext(ctx, &extractedProfiles, q, profileIDs, pag.LastUUID, pag.Limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Profile{}, uuid.Nil, app.NewError(http.StatusNotFound, "user not found", "user not found", err)
		}
		return []domain.Profile{}, uuid.Nil, app.NewError(http.StatusInternalServerError, "unknown error", "postgres internal error", err)
	}

	if len(extractedProfiles) == 0 {
		return []domain.Profile{}, uuid.Nil, nil
	}

	newLastUUID := extractedProfiles[len(extractedProfiles)-1].ID
	if len(extractedProfiles) == 0 {
		newLastUUID = uuid.Nil
	}

	resProfiles := make([]domain.Profile, len(extractedProfiles))

	for i, profile := range extractedProfiles {
		resProfiles[i] = profile.ToProfileDomain()
	}

	return resProfiles, newLastUUID, nil
}
