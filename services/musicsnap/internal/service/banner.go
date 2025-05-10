package service

import (
	"context"
	"fmt"
	"github.com/juju/zaputil/zapctx"
	global "go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"music-snap/pkg/app"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/service/ports"
	"net/http"
)

type BannerService struct {
	cache ports.ProfileCache

	repository ports.BannerRepository
}

func (s BannerService) spanName(funcName string) string {
	return "musicsnap/service." + funcName
}

func (s BannerService) Create(initialCtx context.Context, actor domain.Actor, banner domain.Banner) (int, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(initialCtx, s.spanName("Create"))
	defer span.End()

	ToSpan(&span, actor)

	if !actor.HasRole(domain.AdminRole) {
		return 0, app.NewError(http.StatusForbidden, "user can't create banner",
			fmt.Sprintf("actor do not have %s role", domain.AdminRole), nil)
	}

	bannerInserted, err := s.repository.Create(ctx, banner)
	if err != nil {
		return 0, err
	}
	bannerID := bannerInserted.ID

	return bannerID, nil
}

// TODO rewrite with JOINS
func (s BannerService) GetList(initialCtx context.Context, actor domain.Actor, filter domain.BannerFilter) ([]domain.Banner, error) {
	logger := zapctx.Logger(initialCtx)

	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(initialCtx, s.spanName("GetList"))
	defer span.End()

	ToSpan(&span, actor)

	if !actor.HasOneOfRoles(domain.AdminRole, domain.UserRole) {
		return nil, app.NewError(http.StatusForbidden, "user does not roles to get musicsnap list",
			fmt.Sprintf("actor do not have %s or %s roles", domain.AdminRole, domain.UserRole), nil)
	}

	if filter.Limit == nil {
		filter.Limit = new(int)
		*filter.Limit = 10
	}
	if filter.Offset == nil {
		filter.Offset = new(int)
		*filter.Offset = 0
	}

	banners, err := s.repository.GetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	bannerRes := make([]domain.Banner, 0, len(banners))

	logger.Debug("Span banners", zap.Int("len of banners", len(banners)))
	ToBannersSpan(&span, banners)

	for _, banner := range banners {
		if banner.IsActive || actor.HasRole(domain.AdminRole) {
			bannerRes = append(bannerRes, banner)
		}
	}

	return bannerRes, nil
}

func (s BannerService) Update(initialCtx context.Context, actor domain.Actor, bannerID int, banner domain.Banner) error {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(initialCtx, s.spanName("Update"))
	defer span.End()

	ToSpan(&span, actor)

	if !actor.HasRole(domain.AdminRole) {
		return app.NewError(http.StatusForbidden, "user can't update musicsnap",
			fmt.Sprintf("actor do not have %s role", domain.AdminRole), nil)
	}

	_, err := s.repository.Update(ctx, bannerID, banner)
	if err != nil {
		return err
	}

	return nil
}

func (s BannerService) Delete(initialCtx context.Context, actor domain.Actor, bannerID int) error {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(initialCtx, s.spanName("Delete"))
	defer span.End()

	ToSpan(&span, actor)

	if !actor.HasRole(domain.AdminRole) {
		return app.NewError(http.StatusForbidden, "user can't delete musicsnap",
			fmt.Sprintf("actor do not have %s role", domain.AdminRole), nil)
	}

	err := s.repository.Delete(ctx, bannerID)
	if err != nil {
		return err
	}

	return nil
}

func (s BannerService) Find(initialCtx context.Context, actor domain.Actor, tag, feature int, useLastRevision bool) (domain.CachedBanner, error) {
	tr := global.Tracer(domain.ServiceName)
	ctx, span := tr.Start(initialCtx, s.spanName("Find"))
	defer span.End()

	ToSpan(&span, actor)

	if !actor.HasOneOfRoles(domain.AdminRole, domain.UserRole) {
		return domain.CachedBanner{}, app.NewError(http.StatusForbidden, "user can't find musicsnap",
			fmt.Sprintf("actor do not have %s or %s roles", domain.AdminRole, domain.UserRole), nil)
	}

	cachedBanner, found := s.cache.Get(tag, feature)
	if found {
		return *cachedBanner, nil

	}

	banner, err := s.repository.Find(ctx, tag, feature)
	if err != nil {
		return domain.CachedBanner{}, err
	}

	return domain.CachedBanner{
		ID:       banner.ID,
		Content:  banner.Content,
		IsActive: banner.IsActive,
	}, nil

}
