package storage

import "linksnap/internal/models"

type UrlStore interface {
	GetUrl(code string) (string, error)
	CreateUrl(data *models.NewShortUrl) (*models.ShortUrl, error)
	ListUrlsByUserID(userID string) []*models.ShortUrl
}
