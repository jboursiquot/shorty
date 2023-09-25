package shorty_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jboursiquot/shorty"
)

type stubDDBClient struct {
	av  map[string]types.AttributeValue
	err error
}

func (c *stubDDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return &dynamodb.PutItemOutput{}, c.err
}

func (c *stubDDBClient) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return &dynamodb.GetItemOutput{
		Item: c.av,
	}, c.err
}

func TestDB_Put(t *testing.T) {
	tests := map[string]struct {
		client *stubDDBClient
		su     shorty.ShortenedURL
		err    error
	}{
		"happy path": {
			client: &stubDDBClient{},
			su: shorty.ShortenedURL{
				Key: "bn9Y9rho",
				URL: "https://go.dev",
			},
		},
	}

	ctx := context.Background()

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db := shorty.NewDB(tt.client, "table")
			err := db.Put(ctx, tt.su)
			if err != tt.err {
				t.Errorf("Saver.Save() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestDB_Get(t *testing.T) {
	tests := map[string]struct {
		client *stubDDBClient
		item   shorty.ShortenedURL
		err    error
	}{
		"happy path": {
			client: &stubDDBClient{
				av: map[string]types.AttributeValue{
					"pk":  &types.AttributeValueMemberS{Value: "bn9Y9rho"},
					"url": &types.AttributeValueMemberS{Value: "https://go.dev"},
				},
			},
			item: shorty.ShortenedURL{
				Key: "bn9Y9rho",
				URL: "https://go.dev",
			},
		},
		"not found": {
			client: &stubDDBClient{},
		},
	}

	ctx := context.Background()

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db := shorty.NewDB(tt.client, "table")
			got, err := db.Get(ctx, "key")
			if err != tt.err {
				t.Errorf("Saver.Get() error = %v, want %v", err, tt.err)
				return
			}
			if got == nil {
				return
			}
			if got.Key != tt.item.Key {
				t.Errorf("Saver.Get() got.Key = %v, want %v", got.Key, tt.item.Key)
			}
			if got.URL != tt.item.URL {
				t.Errorf("Saver.Get() got.URL = %v, want %v", got.URL, tt.item.URL)
			}
		})
	}
}
