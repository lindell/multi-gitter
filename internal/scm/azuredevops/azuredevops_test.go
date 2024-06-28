package azuredevops

import (
	"reflect"
	"testing"
)

func TestParseRepositoryReference(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		want    RepositoryReference
		wantErr bool
	}{
		{
			name: "single",
			val:  "my-project/my-repo",
			want: RepositoryReference{
				ProjectName: "my-project",
				Name:        "my-repo",
			},
		},
		{
			name:    "no-project",
			val:     "my-repo",
			wantErr: true,
		},
        {
			name:    "too many parts",
			val:     "my-project/my-repo/more-data",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRepositoryReference(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRepositoryReference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRepositoryReference() = %v, want %v", got, tt.want)
			}
		})
	}
}
