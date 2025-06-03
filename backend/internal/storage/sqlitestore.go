package storage

import (
	"database/sql"
	"fmt"
	"linksnap/internal/models"
)

type SqliteStore struct {
	db *sql.DB
}

func NewSqliteStore(db *sql.DB) *SqliteStore {
	return &SqliteStore{db: db}
}

func (s *SqliteStore) GetUrl(code string) (string, error) {
	var url string
	fmt.Printf("Fetching URL for code: %s\n", code)
	err := s.db.QueryRow(`
		SELECT su.url
		FROM ShortUrl su
		WHERE su.code = ? AND su.deleted_at IS NULL
	`, code).Scan(&url)

	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *SqliteStore) CreateUrl(data *models.NewShortUrl) (*models.ShortUrl, error) {
	var code string

	err := s.db.QueryRow(`
		INSERT INTO ShortUrl (url, code, user_id)
		VALUES (?, ?, ?) RETURNING code
	`, data.Url, data.Code, data.UserId).Scan(&code)

	if err != nil {
		return nil, err
	}

	return &models.ShortUrl{URL: data.Url, Code: code}, nil
}

func (s *SqliteStore) ListUrlsByUserID(userID string) []*models.ShortUrl {
	urls := []*models.ShortUrl{} // slice vacío inicializado

	rows, err := s.db.Query(`
        SELECT id, url, code, created_at
        FROM ShortUrl
        WHERE user_id = ? AND deleted_at IS NULL
        ORDER BY created_at DESC
    `, userID)

	if err != nil {
		return urls // devolver slice vacío con error
	}

	defer rows.Close()

	for rows.Next() {
		var url models.ShortUrl
		url.UserID = userID
		if err := rows.Scan(&url.ID, &url.URL, &url.Code, &url.CreatedAt); err != nil {
			return urls // también devolver slice vacío (o parcial) con error
		}
		urls = append(urls, &url)
	}

	return urls // si todo va bien
}
