package shorty_test

import (
	"testing"

	"github.com/jboursiquot/shorty"
)

func TestShortener_Shorten(t *testing.T) {
	tests := map[string]struct {
		url  string
		want string
	}{
		"go.dev": {
			url:  "https://go.dev",
			want: "bn9Y9rho",
		},
		"pkg.go.dev": {
			url:  "https://pkg.go.dev",
			want: "nnWMuK-X",
		},
	}

	s := shorty.NewShortener()

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := s.Shorten(tt.url); got != tt.want {
				t.Errorf("Shortener.Shorten() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestShortener_Resolve(t *testing.T) {
	tests := map[string]struct {
		key  string
		want string
	}{
		"go.dev": {
			key:  "bn9Y9rho",
			want: "https://go.dev",
		},
		"pkg.go.dev": {
			"nnWMuK-X",
			"https://pkg.go.dev",
		},
	}

	s := shorty.NewShortener()

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			s.Shorten(tt.want)
			if got := s.Resolve(tt.key); got != tt.want {
				t.Errorf("Shortener.Resolve() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
