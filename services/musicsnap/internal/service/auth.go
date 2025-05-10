package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/juju/zaputil/zapctx"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"music-snap/pkg/app"
	pass "music-snap/pkg/password"
	d "music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
)

var _ ports.AuthSvc = &AuthSvc{}

type AuthSvc struct {
	r   ports.UserRepository
	jwt ports.JwtSvc
}

func NewAuthSvc(jwt ports.JwtSvc, userRepository ports.UserRepository) *AuthSvc {
	return &AuthSvc{jwt: jwt,
		r: userRepository}
}

func (s AuthSvc) spanName(funcName string) string {
	//return "musicsnap/service." + reflect.TypeOf(s).NameQuery() + "." + funcName
	return fmt.Sprintf("%s/%s.%s.%s", "musicsnap", "service", "auth", funcName)
}

func (s AuthSvc) EnrichActor(ctx context.Context, actor d.Actor) (d.Actor, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("EnrichActor"))
	defer span.End()

	ToSpan(&span, actor)
	logger := zapctx.Logger(ctx)

	// decode JWT token
	actorFromJWT, err := s.jwt.Parse(actor.Jwt)

	if err != nil {
		return d.Actor{}, app.NewError(http.StatusUnauthorized, "invalid token",
			fmt.Sprintf("invalid token can't parse JWT"), err)
	}

	if (actor.ID != uuid.Nil || actor.ID != uuid.UUID{}) && actor.ID != actorFromJWT.ID {
		// log stolen token
		logger.Warn("stolen token from", zap.String("actorID", actor.ID.String()), zap.String("IDFromJWT", actorFromJWT.ID.String()))
		return d.Actor{}, app.NewError(http.StatusForbidden, "user using alien token",
			fmt.Sprintf("actor %s using stolen token of user %s", actor.ID, actorFromJWT.ID.String()), nil)
	}

	user, err := s.r.GetByID(ctx, actorFromJWT.ID)
	if err != nil {
		return d.Actor{}, err
	}

	askedRoles := d.NewRoles(actor.GetRoles())

	newActor := d.NewActor(user.ID, user.Email, actor.Jwt, user.Nickname, user.Roles.ToSlice())

	newActor.IntersectRoles(askedRoles)

	return newActor, nil
}

func (s AuthSvc) Register(ctx context.Context, actor d.Actor, user d.User, password string) (jwt string, created d.User, err error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Register"))
	defer span.End()

	ToSpan(&span, actor)

	if actor.HasRole(d.UserRole) {
		return "", d.User{},
			app.NewError(http.StatusForbidden, "user already registered",
				fmt.Sprintf("actor  already has role %s", d.UserRole), nil)
	}

	user.PasswordHash, err = pass.HashPassword(password)
	if err != nil {
		return "", d.User{}, err
	}

	if err = s.validForRegistration(ctx, user, password); err != nil {
		return "", d.User{}, err
	}

	user.Roles.Add(d.UserRole)

	userCreated, err := s.r.Create(ctx, user)
	if err != nil {
		return "", d.User{}, err
	}

	// make JWT
	jwt, err = s.jwt.Generate(userCreated)
	if err != nil {
		return "", d.User{}, err
	}

	return jwt, userCreated, nil
}

func (s AuthSvc) validForRegistration(ctx context.Context, user d.User, password string) error {
	err := pass.ValidatePassword(password)
	if err != nil {
		return app.NewError(http.StatusBadRequest, "invalid password",
			fmt.Sprintf("invalid password"), err)
	}

	if user.Email == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for registration",
			fmt.Sprintf("user email is empty"), nil)
	}
	if user.Nickname == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for registration",
			fmt.Sprintf("user name is empty"), nil)
	}
	if user.PasswordHash == "" {
		return app.NewError(http.StatusBadRequest, "invalid user fields for registration",
			fmt.Sprintf("user password is empty"), nil)
	}

	if _, err := s.r.GetByEmail(ctx, user.Email); err == nil {
		return app.NewError(http.StatusBadRequest, "user with same mail already exists",
			fmt.Sprintf("user with this email already exists"), nil)
	}
	//if _, err := s.r.GetByNickname(ctx, user.Nickname); err == nil {
	//	return app.NewError(http.StatusBadRequest, "user with same nickname already exists",
	//		fmt.Sprintf("user with this nickname already exists"), nil)
	//}
	if _, err := s.r.GetByNickname(ctx, user.Nickname); err == nil {
		return app.NewError(http.StatusBadRequest, "user with same mail already exists",
			fmt.Sprintf("user with this email already exists"), nil)
	}
	return nil
}

func (s AuthSvc) Login(ctx context.Context, actor d.Actor, email, password string) (d.User, string, error) {
	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("Login"))
	defer span.End()

	ToSpan(&span, actor)

	user, err := s.r.GetByEmail(ctx, email)
	if err != nil {
		return d.User{}, "", err
	}

	if ok, err := pass.VerifyPassword(password, user.PasswordHash); err != nil || !ok {
		return d.User{}, "", app.NewError(http.StatusUnauthorized, "invalid password",
			fmt.Sprintf("invalid password"), nil)
	}

	jwt, err := s.jwt.Generate(user)
	if err != nil {
		return d.User{}, "", err
	}

	return user, jwt, nil
}

func (s AuthSvc) LogOut(ctx context.Context, actor d.Actor) (d.User, string, error) {
	// TODO

	tr := global.Tracer(d.ServiceName)
	ctx, span := tr.Start(ctx, s.spanName("LogOut"))
	defer span.End()

	ToSpan(&span, actor)

	return d.User{}, "", nil
}

//func (s AuthSvc) Register(ctx context.Context, actor d.Actor, user d.User) (jwt string, err error) {
//	tr := global.Tracer(d.ServiceName)
//	ctx, span := tr.Start(ctx, s.spanName("Create"))
//	defer span.End()
//
//	ToSpan(&span, actor)
//
//	if !actor.HasRole(d.UnregisteredRole) {
//		return "",
//			app.NewError(http.StatusForbidden, "user already registered",
//				fmt.Sprintf("actor do not have %s role", d.UnknownRole), nil)
//	}
//
//	if !(user.Valid() && validForRegistration(user)) {
//		return "",
//			app.NewError(http.StatusForbidden, "invalid user fields for registration",
//				fmt.Sprintf("invalid user fields for registration"), nil)
//	}
//
//	userCreated, err := s.r.Create(ctx, user)
//	if err != nil {
//		return "", err
//	}
//
//	return userCreated, nil
//}
