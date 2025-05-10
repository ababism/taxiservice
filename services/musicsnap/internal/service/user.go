package service

import (
	c "context"
	"fmt"
	"github.com/google/uuid"
	global "go.opentelemetry.io/otel"
	"music-snap/pkg/app"
	"music-snap/pkg/password"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
	"reflect"
)

func (s userSvc) spanName(funcName string) string {
	//return "musicsnap/service." + reflect.TypeOf(s).NameQuery() + "." + funcName
	return fmt.Sprintf("%s/%s.%s.%s", "musicsnap", "service", reflect.TypeOf(s).Name(), funcName)
}

func NewUserSvc(userRepository ports.UserRepository, jwtSvc ports.JwtSvc, cache ports.ProfileCache) ports.UserSvc {
	return userSvc{r: userRepository, c: cache, jwt: jwtSvc}
}

var _ ports.UserSvc = &userSvc{}

type userSvc struct {
	r   ports.UserRepository
	c   ports.ProfileCache
	jwt ports.JwtSvc
}

func (s userSvc) validForCreation(ctx c.Context, user domain.User) error {
	if user.Email == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for creation",
			fmt.Sprintf("validation for creation error user email is empty"), nil)
	}
	if user.Nickname == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for creation",
			fmt.Sprintf("validation for creation error user name is empty"), nil)
	}
	if user.PasswordHash == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for creation",
			fmt.Sprintf("validation for creation error user password is empty"), nil)
	}
	if _, err := s.r.GetByEmail(ctx, user.Email); err == nil {
		return app.NewError(http.StatusBadRequest, "user with same mail already exists",
			fmt.Sprintf("validation for creation error user with this email already exists"), nil)
	}
	if _, err := s.r.GetByNickname(ctx, user.Nickname); err == nil {
		return app.NewError(http.StatusBadRequest, "user with same nickname already exists",
			fmt.Sprintf("validation for creation error user with this nickname already exists"), nil)
	}
	if !user.Valid() {
		return app.NewError(http.StatusBadRequest, "invalid user fields for creation",
			fmt.Sprintf("validation for creation error invalid user fields for creation"), nil)
	}
	return nil
}

func (s userSvc) Create(ctx c.Context, actor domain.Actor, user domain.User, pass string) (domain.User, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Create"))
	defer span.End()
	ToSpan(&span, actor)

	if !actor.HasRole(domain.AdminRole) {
		return domain.User{},
			app.NewError(http.StatusForbidden, "user can't create user",
				"actor do not have admin role", nil)
	}

	var err error
	user.PasswordHash, err = password.HashPassword(pass)
	if err != nil {
		return domain.User{}, err
	}

	user.Roles.Add(domain.UserRole)
	err = s.validForCreation(ctx, user)
	if err != nil {
		return domain.User{},
			app.NewError(http.StatusForbidden, "invalid user fields for registration",
				"invalid user fields for registration", err)
	}

	userCreated, err := s.r.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	return userCreated, nil
}

func (s userSvc) Get(ctx c.Context, actor domain.Actor, userID uuid.UUID) (domain.User, error) {
	if !actor.HasRole(domain.AdminRole) {
		if actor.ID != userID {
			return domain.User{},
				app.NewError(http.StatusForbidden, "user can't get other user private data",
					"actor do not have admin role", nil)
		}
	}
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Get"))
	defer span.End()
	ToSpan(&span, actor)

	user, err := s.r.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, app.NewError(http.StatusBadRequest, "user not found",
			"can't get user by actor id", err)
	}

	return user, nil
}

func (s userSvc) validForUpdate(ctx c.Context, user domain.User) error {
	if user.Email == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for update",
			fmt.Sprintf("user email is empty"), nil)
	}
	if user.Nickname == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for update",
			fmt.Sprintf("user name is empty"), nil)
	}
	if user.PasswordHash == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for update",
			fmt.Sprintf("user password is empty"), nil)
	}
	if grabbedUser, err := s.r.GetByEmail(ctx, user.Email); err == nil && grabbedUser.ID != user.ID {
		return app.NewError(http.StatusBadRequest, "user with same mail already exists",
			fmt.Sprintf("user with this email already exists"), nil)
	}
	if grabbedUser, err := s.r.GetByNickname(ctx, user.Nickname); err == nil && user.ID != grabbedUser.ID {
		return app.NewError(http.StatusBadRequest, "user with same nickname already exists",
			fmt.Sprintf("user with this nickname already exists"), nil)
	}
	if !user.Valid() {
		return app.NewError(http.StatusBadRequest, "invalid user fields for update",
			fmt.Sprintf("invalid user fields for update"), nil)
	}
	return nil
}

func (s userSvc) Update(ctx c.Context, actor domain.Actor, user domain.User, pass string) (domain.User, error) {
	if !actor.HasRole(domain.AdminRole) && actor.ID != user.ID {
		return domain.User{},
			app.NewError(http.StatusForbidden, "user can't update other user",
				"actor do not have admin role or can't update other user", nil)
	}
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Update"))
	defer span.End()
	ToSpan(&span, actor)

	var err error
	user.PasswordHash, err = password.HashPassword(pass)
	if err != nil {
		return domain.User{}, err
	}

	if err := s.validForUpdate(ctx, user); err != nil {
		return domain.User{},
			app.NewError(http.StatusForbidden, "invalid user fields for update",
				"invalid user fields for update", err)
	}

	userUpdated, err := s.r.UpdateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return userUpdated, nil
}

func (s userSvc) GetProfile(ctx c.Context, actor domain.Actor, userID uuid.UUID) (domain.Profile, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Get"))
	defer span.End()

	ToSpan(&span, actor)

	user, err := s.r.GetByID(ctx, userID)
	if err != nil {
		return domain.Profile{}, app.NewError(http.StatusBadRequest, "user not found",
			"can't get user by actor id", err)
	}

	return user.Profile, nil
}

func (s userSvc) GetProfilesList(ctx c.Context, actor domain.Actor, nickNameQuery string, pagination domain.UUIDPagination) ([]domain.Profile, domain.UUIDPagination, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("GetProfilesList"))
	defer span.End()

	ToSpan(&span, actor)

	users, newLastUUID, err := s.r.GetList(ctx, nickNameQuery, pagination)
	pagination.LastUUID = newLastUUID
	if err != nil {
		return []domain.Profile{}, pagination, app.NewError(http.StatusBadRequest, "users search error",
			"can't get users by search query", err)
	}

	profiles := make([]domain.Profile, len(users))
	for i, user := range users {
		profiles[i] = user.Profile
	}
	return profiles, pagination, nil
}

func (s userSvc) GetUserList(ctx c.Context, actor domain.Actor, nickNameQuery string, pagination domain.UUIDPagination) ([]domain.User, domain.UUIDPagination, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("GetUserList"))
	defer span.End()

	ToSpan(&span, actor)

	users, newLastUUID, err := s.r.GetList(ctx, nickNameQuery, pagination)
	pagination.LastUUID = newLastUUID
	if err != nil {
		return []domain.User{}, pagination, app.NewError(http.StatusBadRequest, "users search error",
			"can't get users by query", err)
	}

	return users, pagination, nil
}

//
//func (s userSvc) UpdatePr(ctx c.Context, actor domain.Actor, profile domain.Profile) (domain.Profile, error) {
//	tr := global.Tracer(domain.ServiceName)
//	ctx, span := tr.Start(ctx, s.spanName("Create"))
//	defer span.End()
//
//	ToSpan(&span, actor)
//
//	if !actor.HasRole(domain.RegisteredRole) {
//		return domain.Profile{},
//			app.NewError(http.StatusForbidden, "user not registered",
//				"actor do not have registered role", nil)
//	}
//
//	user, err := s.r.GetByID(ctx, actor.ID)
//	if err != nil {
//		return domain.Profile{}, app.NewError(http.StatusBadRequest, "user not found",
//			"can't get user by actor id", err)
//	}
//
//	user.Profile = profile
//
//	userCreated, err := s.r.UpdateUser(ctx, user)
//	if err != nil {
//		return domain.Profile{}, err
//	}
//
//	return userCreated.Profile, nil
//}

func (s userSvc) UpdateProfile(ctx c.Context, actor domain.Actor, profile domain.Profile) (domain.Profile, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("UpdateProfile"))
	defer span.End()

	ToSpan(&span, actor)

	if !actor.HasRole(domain.UserRole) {
		return domain.Profile{},
			app.NewError(http.StatusForbidden, "user not registered",
				"actor do not have registered role", nil)
	}

	user, err := s.r.GetByID(ctx, actor.ID)
	if err != nil {
		return domain.Profile{}, app.NewError(http.StatusBadRequest, "user not found",
			"can't get user by actor id", err)
	}

	user.Profile = profile

	userCreated, err := s.r.UpdateUser(ctx, user)
	if err != nil {
		return domain.Profile{}, err
	}

	return userCreated.Profile, nil
}

//func (s userSvc) Create(ctx c.Context, actor domain.Actor, user domain.User) (domain.User, error) {
//	tr := global.Tracer(domain.ServiceName)
//	ctx, span := tr.Start(ctx, s.spanName("Create"))
//	defer span.End()
//
//	ToSpan(&span, actor)
//
//	if !actor.HasRole(domain.UnregisteredRole) {
//		return domain.User{},
//			app.NewError(http.StatusForbidden, "user already registered",
//				fmt.Sprintf("actor do not have %s role", domain.UnknownRole), nil)
//	}

//if !(user.Valid() && validForRegistration(user)) {
//	return domain.User{},
//		app.NewError(http.StatusForbidden, "invalid user fields for registration",
//			fmt.Sprintf("invalid user fields for registration"), nil)
//}
//
//	userCreated, err := s.r.Create(ctx, user)
//	if err != nil {
//		return domain.User{}, err
//	}
//
//	return userCreated, nil
//}
//

//func (s userSvc) Create(initialCtx c.Context, actor domain.Actor, banner domain.Banner) (int, error) {
//	tr := global.Tracer(domain.ServiceName)
//	ctx, span := tr.Start(initialCtx, s.spanName("Create"))
//	defer span.End()
//
//	ToSpan(&span, actor)
//
//	if !actor.HasRole(domain.AdminRole) {
//		return 0, app.NewError(http.StatusForbidden, "user can't create banner",
//			fmt.Sprintf("actor do not have %s role", domain.AdminRole), nil)
//	}
//
//	bannerInserted, err := s.r.Create(ctx, banner)
//	if err != nil {
//		return 0, err
//	}
//	bannerID := bannerInserted.ID
//
//	return bannerID, nil
//}
