package service

import (
	"context"
	"fmt"
	"linksnap/internal/models"
	"linksnap/internal/storage"
	"linksnap/internal/util"
	"time"

	"github.com/redis/go-redis/v9"
)

type ShortUrlService struct {
	store storage.UrlStore
	cache redis.Client
}

func NewShortUrlService(store storage.UrlStore, redis redis.Client) *ShortUrlService {
	return &ShortUrlService{
		store: store,
		cache: redis,
	}
}

func (s *ShortUrlService) ListUrlsByUserID(userID string) []*models.ShortUrl {
	return s.store.ListUrlsByUserID(userID)
}

func (s *ShortUrlService) CreateUrl(userId string, url string) (*models.ShortUrl, error) {
	urls := s.store.ListUrlsByUserID(userId)

	if len(urls) >= 10 {
		return nil, fmt.Errorf("Maximo de urls alcanzado")
	}

	code, err := util.GenerateRandomCode(10)
	if err != nil {
		return nil, err
	}

	return s.store.CreateUrl(&models.NewShortUrl{
		UserId: userId,
		Url:    url,
		Code:   code,
	})
}

func (s *ShortUrlService) DeleteUrl(userId string, id string) (*models.ShortUrl, error) {
	return s.store.DeleteUrl(userId, id)
}

func (s *ShortUrlService) Resolve(ctx context.Context, code string) (string, error) {

	url, err := s.cache.Get(ctx, code).Result()
	if err == redis.Nil {
		url, err = s.store.GetUrl(code)
		if err != nil {
			return "", err
		}

		if url == "" {
			return "", fmt.Errorf("No existe la url consultada")
		}

		s.cache.Set(ctx, code, url, time.Hour*4)

		return url, nil
	}

	if err != nil {
		return "", err
	}

	return url, nil
}
