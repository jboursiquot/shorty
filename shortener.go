package shorty

import (
	"crypto/sha256"
	"encoding/base64"
	"sync"
)

type Shortener struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewShortener() *Shortener {
	return &Shortener{
		urls: make(map[string]string),
	}
}

func (s *Shortener) Shorten(url string) string {
	if u := s.Resolve(url); u != "" {
		return u
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	hash := sha256.Sum256([]byte(url))
	key := base64.URLEncoding.EncodeToString(hash[:])[:8]
	s.urls[key] = url
	return key
}

func (s *Shortener) Resolve(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if existing, ok := s.urls[key]; ok {
		return existing
	}
	return ""
}
