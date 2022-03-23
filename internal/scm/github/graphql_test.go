package github

import (
	"fmt"
	"testing"
)

func Test_graphQLEndpoint(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{url: "https://github.detsbihcs.io/api/v3", want: "https://github.detsbihcs.io/api/graphql"},
		{url: "https://github.detsbihcs.io/api/v3/", want: "https://github.detsbihcs.io/api/graphql"},
		{url: "https://github.detsbihcs.io/api/", want: "https://github.detsbihcs.io/api/graphql"},
		{url: "https://github.detsbihcs.io/", want: "https://github.detsbihcs.io/api/graphql"},
		{url: "https://api.github.detsbihcs.io/", want: "https://api.github.detsbihcs.io/graphql"},
		{url: "https://more.api.github.detsbihcs.io/", want: "https://more.api.github.detsbihcs.io/graphql"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("url: %s", tt.url), func(t *testing.T) {
			if got, _ := graphQLEndpoint(tt.url); got != tt.want {
				t.Errorf("graphQLEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
