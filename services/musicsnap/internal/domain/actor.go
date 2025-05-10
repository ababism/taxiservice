package domain

import (
	"github.com/google/uuid"
	"music-snap/pkg/stringset"
)

// Actor represents a user of system with data taken from request
type Actor struct {
	ID       uuid.UUID
	Mail     string
	Jwt      string
	Nickname string
	// текущие роли
	// roles will be slice of strings in API layer
	roles stringset.Set
}

func NewActor(id uuid.UUID, mail, jwt, nickname string, roles []string) Actor {
	rolesSet := stringset.New(roles...)
	//rolesSet.AddItems(roles)
	return Actor{
		ID:       id,
		Mail:     mail,
		Jwt:      jwt,
		Nickname: nickname,
		roles:    rolesSet,
	}
}

// TODO remove legacy
func NewActorFromRoles(roles []string) Actor {
	a := Actor{
		//ID: ID,
		roles: stringset.New()}
	a.initRoles(roles)
	return a
}

func (a *Actor) HasRole(role string) bool {
	return a.roles.Contains(role)
}

func (a *Actor) AreRolesEmpty() bool {
	return a.roles.Size() == 0
}

func (a *Actor) HasOneOfRoles(roles ...string) bool {
	if roles == nil || len(roles) == 0 {
		return true
	}

	for _, role := range roles {
		if a.roles.Contains(role) {
			return true
		}
	}
	return false
}

func (a *Actor) HasAllRoles(roles ...string) bool {
	if roles == nil || len(roles) == 0 {
		return true
	}

	for _, role := range roles {
		if !a.roles.Contains(role) {
			return false
		}
	}
	return true
}

func (a *Actor) initRoles(roles []string) {
	a.roles.AddItems(roles)
}

func (a *Actor) GetRoles() []string {
	// TODO - does it necessary to return make([]string, 0)
	if a == nil || a.roles == nil {
		return make([]string, 0)
	}
	copyRoles := make([]string, 0, a.roles.Size())
	for r, _ := range a.roles {
		copyRoles = append(copyRoles, r)
	}
	return copyRoles
}
func (a *Actor) AddRole(role string) {
	if role == "" {
		return
	}
	a.roles.Add(role)
}

func (a *Actor) AddRoles(roles []string) {
	if roles == nil || len(roles) == 0 {
		return
	}
	for _, role := range roles {
		a.roles.Add(role)
	}
}
func (a *Actor) IntersectRoles(other Roles) {
	if a == nil || other.Empty() {
		return
	}
	a.roles = a.roles.Intersect(other.getRolesSet())
}

func (c *BannerFilter) Validate() error {
	//app.NewError(http.StatusBadRequest, " ", " ", nil)
	return nil
}
