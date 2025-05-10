package stringset

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	studentKey = "studentKey"
	teacherKey = "teacherKey"
	adminKey   = "adminKey"
	unknownKey = "unknownKey"
)

func TestSet(t *testing.T) {
	t.Parallel()

	t.Run("Test adding and removing elements", func(t *testing.T) {
		t.Parallel()

		set := make(Set)

		set.Add(studentKey)
		set.Add(teacherKey)
		set.Add(adminKey)

		assert.True(t, set.Contains(studentKey))
		assert.True(t, set.Contains(teacherKey))
		assert.True(t, set.Contains(adminKey))
		assert.False(t, set.Contains(unknownKey))

		set.Remove(teacherKey)

		assert.False(t, set.Contains(teacherKey))
		assert.Equal(t, 2, set.Size())
	})

	t.Run("Test getting items", func(t *testing.T) {
		t.Parallel()

		set := make(Set)
		set.Add(studentKey)
		set.Add(teacherKey)
		set.Add(adminKey)

		items := set.Items()

		assert.Len(t, items, 3)
		assert.Contains(t, items, studentKey)
		assert.Contains(t, items, teacherKey)
		assert.Contains(t, items, adminKey)
	})

	t.Run("Test nil set", func(t *testing.T) {
		t.Parallel()

		set := make(Set)
		set = nil

		assert.False(t, set.Contains(studentKey))
	})
}

func TestSet_AddItems(t *testing.T) {
	t.Parallel()

	t.Run("Test adding items", func(t *testing.T) {
		t.Parallel()

		set := make(Set)

		set.AddItems([]string{studentKey, teacherKey, adminKey})

		assert.True(t, set.Contains(studentKey))
		assert.True(t, set.Contains(teacherKey))
		assert.True(t, set.Contains(adminKey))
		assert.Equal(t, 3, set.Size())
	})

	t.Run("Test adding items to an existing set", func(t *testing.T) {
		t.Parallel()

		set := make(Set)
		set.Add(studentKey)

		set.AddItems([]string{teacherKey, adminKey})

		assert.True(t, set.Contains(studentKey))
		assert.True(t, set.Contains(teacherKey))
		assert.True(t, set.Contains(adminKey))
		assert.Equal(t, 3, set.Size())
	})
}

func TestSet_Intersect(t *testing.T) {
	tests := []struct {
		name     string
		receiver Set
		other    Set
		want     Set
	}{
		{
			name:     "empty sets",
			receiver: New(),
			other:    New(),
			want:     New(),
		},
		{
			name:     "non-empty with empty",
			receiver: New("a", "b", "c"),
			other:    New(),
			want:     New(),
		},
		{
			name:     "empty with non-empty",
			receiver: New(),
			other:    New("a", "b", "c"),
			want:     New(),
		},
		{
			name:     "some intersection",
			receiver: New("a", "b", "c"),
			other:    New("b", "c", "d"),
			want:     New("b", "c"),
		},
		{
			name:     "full intersection",
			receiver: New("a", "b"),
			other:    New("a", "b"),
			want:     New("a", "b"),
		},
		{
			name:     "no intersection",
			receiver: New("a", "b"),
			other:    New("c", "d"),
			want:     New(),
		},
		{
			name:     "nil receiver",
			receiver: nil,
			other:    New("a", "b"),
			want:     New(),
		},
		{
			name:     "nil other set",
			receiver: New("a", "b"),
			other:    nil,
			want:     New(),
		},
		{
			name:     "both nil",
			receiver: nil,
			other:    nil,
			want:     New(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange - test case setup is already done in the struct

			// Act
			got := tt.receiver.Intersect(tt.other)

			// Assert
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Intersect() mismatch (-want +got):\n%s", diff)
			}

			// Additional assertion to ensure original set wasn't modified
			if tt.receiver != nil && tt.other != nil && len(tt.receiver) != len(New(tt.receiver.ToSlice()...)) {
				t.Error("Intersect() modified the original set, want immutable operation")
			}
		})
	}
}

// Helper method to convert set to slice for verification
func (s Set) ToSlice() []string {
	var items []string
	for item := range s {
		items = append(items, item)
	}
	return items
}
