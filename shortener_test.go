package shorty_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jboursiquot/shorty"
)

func TestShortener_Shorten(t *testing.T) {
	tests := map[string]struct {
		url    string
		want   string
		client *stubDDBClient
	}{
		"go.dev": {
			url:    "https://go.dev",
			want:   "bn9Y9rho",
			client: &stubDDBClient{},
		},
		"pkg.go.dev": {
			url:    "https://pkg.go.dev",
			want:   "nnWMuK-X",
			client: &stubDDBClient{},
		},
	}

	ctx := context.Background()
	saver := shorty.NewDB(&stubDDBClient{}, "table")
	s := shorty.NewShortener(saver)

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := s.Shorten(ctx, tt.url); got != tt.want {
				t.Errorf("Shortener.Shorten() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func TestShortener_Resolve(t *testing.T) {
	tests := map[string]struct {
		key    string
		want   string
		client *stubDDBClient
	}{
		"go.dev": {
			key:  "bn9Y9rho",
			want: "https://go.dev",
			client: &stubDDBClient{
				av: map[string]types.AttributeValue{
					"pk":  &types.AttributeValueMemberS{Value: "bn9Y9rho"},
					"url": &types.AttributeValueMemberS{Value: "https://go.dev"},
				},
			},
		},
		"pkg.go.dev": {
			"nnWMuK-X",
			"https://pkg.go.dev",
			&stubDDBClient{
				av: map[string]types.AttributeValue{
					"pk":  &types.AttributeValueMemberS{Value: "nnWMuK-X"},
					"url": &types.AttributeValueMemberS{Value: "https://pkg.go.dev"},
				},
			},
		},
	}
	ctx := context.Background()

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			saver := shorty.NewDB(tt.client, "table")
			s := shorty.NewShortener(saver)
			if got := s.Resolve(ctx, tt.key); got != tt.want {
				t.Errorf("Shortener.Resolve() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
