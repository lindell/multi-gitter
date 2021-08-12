package gitlab

import (
	"reflect"
	"testing"
)

func TestParseProjectReference(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    ProjectReference
		wantErr bool
	}{
		{
			name: "single",
			val:  "my-group/my-project",
			want: ProjectReference{
				OwnerName: "my-group",
				Name:      "my-project",
			},
		},
		{
			name: "subgroup",
			val:  "my-group/sub-group/my-project",
			want: ProjectReference{
				OwnerName: "my-group/sub-group",
				Name:      "my-project",
			},
		},
		{
			name: "two subgroups",
			val:  "my-group/sub-group1/sub-group2/my-project",
			want: ProjectReference{
				OwnerName: "my-group/sub-group1/sub-group2",
				Name:      "my-project",
			},
		},
		{
			name:    "no-group",
			val:     "my-project",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProjectReference(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProjectReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseProjectReference() = %v, want %v", got, tt.want)
			}
		})
	}
}
