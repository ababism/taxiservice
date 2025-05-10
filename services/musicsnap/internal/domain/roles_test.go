package domain

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRoles(t *testing.T) {
	t.Parallel()
	// Arrange a new roles for testing
	const (
		role1        = "role1"
		role2        = "role2"
		missingRole1 = "missingRole1"
	)

	roles := NewRoles([]string{role1, role2})

	t.Run("HasOneOf", func(t *testing.T) {
		t.Parallel()

		assert.True(t, roles.HasOneOf(role1))
		assert.True(t, roles.HasOneOf(role2))
		assert.True(t, roles.HasOneOf(role1, missingRole1))
		assert.False(t, roles.HasOneOf(missingRole1, missingRole1))
		assert.False(t, roles.HasOneOf(missingRole1))

	})

	t.Run("HasOneOfRolesEmptyRoles", func(t *testing.T) {
		t.Parallel()

		// Test bad scenario
		assert.True(t, roles.HasOneOf())
	})

	t.Run("HasOneOfRolesEmptyRoles", func(t *testing.T) {
		t.Parallel()

		emptyRoles := NewRoles(make([]string, 0))

		assert.False(t, emptyRoles.HasOneOf(role1))
		assert.False(t, emptyRoles.HasOneOf(missingRole1))
	})

	t.Run("HasAll", func(t *testing.T) {
		t.Parallel()

		assert.True(t, roles.HasAll(role1, role2))
		assert.False(t, roles.HasAll(role1, missingRole1))
		assert.False(t, roles.HasAll(missingRole1))
	})

	t.Run("InitRoles", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, 2, roles.roles.Size(), "unexpected number of roles")
		assert.True(t, roles.roles.Contains(role1))
		assert.True(t, roles.roles.Contains(role2))
	})
}

func TestHasOneOf_WithRole_ReturnsTrue(t *testing.T) {
	t.Parallel()

	roles := NewRoles([]string{"role1", "role2"})

	assert.True(t, roles.HasOneOf("role1"))
}

// Может так лучше?
func TestHasOneOf_WithoutRole_ReturnsFalse(t *testing.T) {
	t.Parallel()

	roles := NewRoles([]string{"role1", "role2"})

	assert.False(t, roles.HasOneOf("role3"))
}

func TestHasOneOf_WithEmptyRoles_ReturnsTrue(t *testing.T) {
	t.Parallel()

	roles := NewRoles([]string{"role1", "role2"})

	assert.True(t, roles.HasOneOf())
}

func TestHasOneOf_WithNilRoles_ReturnsTrue(t *testing.T) {
	t.Parallel()

	roles := NewRoles([]string{"role1", "role2"})

	assert.True(t, roles.HasOneOf(nil...))
}

func TestGet_ReturnsCorrectRoles(t *testing.T) {
	t.Parallel()
	arrangedRoles := []string{"role1", "role2"}

	roles := NewRoles(arrangedRoles)

	returnedRoles := roles.ToSlice()

	assert.ElementsMatch(t, arrangedRoles, returnedRoles)
}

func TestToSlice_NoRoles_ReturnsEmptySlice(t *testing.T) {
	t.Parallel()

	roles := NewRoles(nil)

	returnedRoles := roles.ToSlice()

	assert.Empty(t, returnedRoles)
}
