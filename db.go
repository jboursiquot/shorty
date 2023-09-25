package shorty

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DDBClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type DB struct {
	client DDBClient
	table  string
}

func NewDB(client DDBClient, table string) *DB {
	return &DB{
		client: client,
		table:  table,
	}
}

type ShortenedURL struct {
	Key string `json:"key" dynamodbav:"pk"`
	URL string `json:"url" dynamodbav:"url"`
}

func (s *DB) Put(ctx context.Context, item ShortenedURL) error {
	input := &dynamodb.PutItemInput{
		TableName: &s.table,
		Item: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: item.Key},
			"url": &types.AttributeValueMemberS{
				Value: item.URL,
			},
		},
	}

	_, err := s.client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}

	return nil
}

func (s *DB) Get(ctx context.Context, key string) (*ShortenedURL, error) {
	input := &dynamodb.GetItemInput{
		TableName: &s.table,
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: key},
		},
	}

	result, err := s.client.GetItem(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %v", err)
	}

	if result.Item == nil {
		slog.Info("GetItemOutput", "result", "item is nil")
		return nil, nil
	}

	var item ShortenedURL
	if err := attributevalue.UnmarshalMap(result.Item, &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %v", err)
	}

	return &item, nil
}
