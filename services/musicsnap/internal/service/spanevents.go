package service

import (
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"music-snap/services/musicsnap/internal/domain"
	"music-snap/services/musicsnap/internal/domain/keys"
)

func ToSpan(span *trace.Span, a domain.Actor) {
	if span == nil {
		return
	}
	//(*span).SetAttributes(attribute.String(keys.ActorIDAttributeKey, a.ID.String()))
	(*span).SetAttributes(attribute.StringSlice(keys.ActorRolesAttributeKey, a.GetRoles()))

	(*span).AddEvent("actor logged in for service method")
}

func ToBannersSpan(span *trace.Span, bs []domain.Banner) {
	if span == nil {
		return
	}
	//(*span).SetAttributes(attribute.String(keys.ActorIDAttributeKey, a.ID.String()))
	//(*span).SetAttributes(attribute.StringSlice(keys.ActorRolesAttributeKey, a.GetRoles()))
	bannerStrings := make([]string, len(bs))
	for i, v := range bs {
		bannerStrings[i] = fmt.Sprintf("{id: %d}", v.ID)
	}
	(*span).AddEvent("test event")
	(*span).AddEvent("banners from repository: ", trace.WithAttributes(attribute.StringSlice("banners", bannerStrings)))
}
