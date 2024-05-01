package scm

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name        string
		s1          []int
		s2          []int
		wantAdded   []int
		wantRemoved []int
	}{
		{
			name:        "same",
			s1:          []int{1, 2, 3},
			s2:          []int{1, 2, 3},
			wantAdded:   nil,
			wantRemoved: nil,
		},
		{
			name:        "empty s2",
			s1:          []int{1, 2, 3},
			s2:          []int{},
			wantAdded:   nil,
			wantRemoved: []int{1, 2, 3},
		},
		{
			name:        "empty s1",
			s1:          []int{},
			s2:          []int{1, 2, 3},
			wantAdded:   []int{1, 2, 3},
			wantRemoved: nil,
		},
		{
			name:        "some overlap",
			s1:          []int{1, 2, 3},
			s2:          []int{3, 4, 5},
			wantAdded:   []int{4, 5},
			wantRemoved: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdded, gotRemoved := Diff(tt.s1, tt.s2)
			if !reflect.DeepEqual(gotAdded, tt.wantAdded) {
				t.Errorf("Diff() gotAdded = %v, want %v", gotAdded, tt.wantAdded)
			}
			if !reflect.DeepEqual(gotRemoved, tt.wantRemoved) {
				t.Errorf("Diff() gotRemoved = %v, want %v", gotRemoved, tt.wantRemoved)
			}
		})
	}
}
