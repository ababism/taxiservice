package musicsnap

import (
	"github.com/google/uuid"
	qb "music-snap/pkg/querybuilder"
	"testing"
	//"github.com/stretchr/testify/require"
)

type ReviewFilter struct {
	UserID  *uuid.UUID
	PieceID *string

	Rating    *int // 1-10
	Moderated *bool
	Published *bool

	IncludeProfiles bool
	OrderByRating   *bool
	OrderAsc        bool
}

// write unit test if uuid.ni; == uuid.UUID{}
func TestQB(t *testing.T) {

	userID := uuid.New()
	pieceID := "someSpotifyID"
	rating := 6
	moderated := true
	published := true

	_ = userID
	_ = pieceID
	_ = rating
	_ = moderated
	_ = published

	orderByRating := false
	filter := ReviewFilter{
		//UserID:          &userID,
		//PieceID: &pieceID,
		//Rating:          &rating,
		//Moderated: &moderated,
		//Published:       &published,
		IncludeProfiles: true,
		OrderByRating:   &orderByRating,
		OrderAsc:        true,
	}

	orderByField := "created_at"
	if filter.OrderByRating != nil && *filter.OrderByRating {
		orderByField = "rating"
	}

	qBuild := qb.NewNamed().
		Q("SELECT * FROM reviews").
		StartOpt().
		Q("JOIN").Table("profiles").ON().Q("reviews.user_id = profiles.user_id").
		EndOptIf(func() bool {
			return filter.IncludeProfiles
		}).
		WhereOptPart().
		CompConnectorOpt("user_id", qb.EQ(), "user_id", filter.UserID, qb.AND()).
		CompConnectorOpt("piece_id", qb.EQ(), "piece_id", filter.PieceID, qb.AND()).
		CompConnectorOpt("rating", qb.GET(), "rating", filter.Rating, qb.AND()).
		CompConnectorOpt("moderated", qb.EQ(), "moderated", filter.Moderated, qb.AND()).
		CompConnectorOpt("published", qb.EQ(), "published", filter.Published, qb.AND()).
		EndWhereOpt().
		OrderBy(orderByField, filter.OrderAsc).
		Limit("last_id", 10)

	q, args := qBuild.Build()

	t.Log(q)
	t.Log(args)

	//filter := struct {
	//	Published *bool
	//	LastID    *int
	//	UserID    *uuid.UUID
	//}{
	//	Published: nil,
	//	LastID:    nil,
	//	UserID:    &uuidTest,
	//}
	//qb := querybuilder.NewNamed().Q(
	//	"SELECT * FROM reviews").
	//	WhereOptPart().CompOpt("published", querybuilder.EQ(), "published", filter.Published).
	//	CompOpt("user_id", querybuilder.EQ(), "user_id", filter.UserID).
	//	Limit(nil, 10)
	//
	//q, args := qb.Build()
	//
	//t.Log(q)
	//t.Log(args)
	//
	//qb = querybuilder.NewNamed().Q(
	//	"SELECT * FROM reviews").
	//	WhereOptPart().Col("published").EQ().ArgOpt("published", filter.Published).
	//	CompOpt("user_id", querybuilder.EQ(), "user_id", filter.UserID).Limit(nil, 10)
	//
	//q, args = qb.Build()
	//
	//t.Log(q)
	//t.Log(args)
}
