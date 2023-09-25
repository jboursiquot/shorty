package shorty

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
)

type Shortener struct {
	db *DB
}

func NewShortener(db *DB) *Shortener {
	return &Shortener{
		db: db,
	}
}

func (s *Shortener) Shorten(ctx context.Context, url string) string {
	hash := sha256.Sum256([]byte(url))
	key := base64.URLEncoding.EncodeToString(hash[:])[:8]

	su := ShortenedURL{
		Key: key,
		URL: url,
	}

	if err := s.db.Put(ctx, su); err != nil {
		return ""
	}

	return key
}

func (s *Shortener) Resolve(ctx context.Context, key string) string {
	if su, err := s.db.Get(ctx, key); err == nil && su != nil {
		return su.URL
	}

	return ""
}
