package domain

import (
	"music-snap/pkg/stringset"
)

const (
	// Role to avoid using empty data structures
	NoRole = "no_role"
	//RegisteredRole = "registered"
	UserRole = "user"
	//UnregisteredRole = "unregistered"
	AdminRole     = "admin"
	ModeratorRole = "moderator"
)

// Roles represents roles of user
type Roles struct {
	// текущие роли
	roles stringset.Set
}

func NewRoles(roles []string) Roles {
	a := Roles{
		roles: stringset.New()}
	a.roles.AddItems(roles)
	return a
}

func (rs *Roles) Has(role string) bool {
	return rs.roles.Contains(role)
}

func (rs *Roles) Empty() bool {
	if rs == nil || rs.roles == nil {
		return true
	}
	return rs.roles.Size() == 0
}

func (rs *Roles) HasOneOf(roles ...string) bool {
	if roles == nil || len(roles) == 0 {
		return true
	}

	for _, role := range roles {
		if rs.roles.Contains(role) {
			return true
		}
	}
	return false
}

func (rs *Roles) HasAll(roles ...string) bool {
	if roles == nil || len(roles) == 0 {
		return true
	}

	for _, role := range roles {
		if !rs.roles.Contains(role) {
			return false
		}
	}
	return true
}

func (rs *Roles) ToSlice() []string {
	// TODO - does it necessary to return make([]string, 0)
	if rs == nil || rs.roles == nil {
		return make([]string, 0)
	}
	copyRoles := make([]string, 0, rs.roles.Size())
	for r, _ := range rs.roles {
		copyRoles = append(copyRoles, r)
	}
	return copyRoles
}

func (rs *Roles) Add(role string) {
	if role == "" {
		return
	}
	rs.roles.Add(role)
}

func (rs *Roles) AddRoles(roles []string) {
	// generate code
	if roles == nil || len(roles) == 0 {
		return
	}
	for _, role := range roles {
		rs.roles.Add(role)
	}
}

func (rs *Roles) Intersect(other Roles) {
	if rs == nil || other.Empty() {
		return
	}
	if rs.Empty() {
		rs.roles = stringset.New()
		return
	}

	rs.roles.Intersect(other.roles)
	if rs.roles.Size() == 0 {
		rs.roles = stringset.New()
	}
}

func (rs *Roles) getRolesSet() stringset.Set {
	if rs == nil || rs.roles == nil {
		return stringset.New()
	}
	return rs.roles
}
