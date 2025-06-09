package storage

import (
	"database/sql"
	"fmt"
	"linksnap/internal/models"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

// Obtener URL original a partir del código
func (s *PostgresStore) GetUrl(code string) (string, error) {
	var url string
	fmt.Printf("Fetching URL for code: %s\n", code)
	err := s.db.QueryRow(`
		SELECT su.url
		FROM shorturl su
		WHERE su.code = $1 AND su.deleted_at IS NULL
	`, code).Scan(&url)

	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return "", err
	}

	return url, nil
}

// Crear nueva URL acortada
func (s *PostgresStore) CreateUrl(data *models.NewShortUrl) (*models.ShortUrl, error) {
	var code string

	err := s.db.QueryRow(`
		INSERT INTO shorturl (url, code, user_id)
		VALUES ($1, $2, $3)
		RETURNING code
	`, data.Url, data.Code, data.UserId).Scan(&code)

	if err != nil {
		fmt.Println("Error creating URL:", err)
		return nil, err
	}

	return &models.ShortUrl{URL: data.Url, Code: code}, nil
}

// Listar URLs por ID de usuario
func (s *PostgresStore) ListUrlsByUserID(userID string) []*models.ShortUrl {
	urls := []*models.ShortUrl{}

	rows, err := s.db.Query(`
		SELECT id, url, code, created_at
		FROM shorturl
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		fmt.Println("Error fetching URLs:", err)
		return urls
	}

	defer rows.Close()

	for rows.Next() {
		var url models.ShortUrl
		url.UserID = userID
		if err := rows.Scan(&url.ID, &url.URL, &url.Code, &url.CreatedAt); err != nil {
			return urls
		}
		urls = append(urls, &url)
	}

	return urls
}

func (s *PostgresStore) DeleteUrl(userId string, id string) (*models.ShortUrl, error) {
	var url models.ShortUrl

	err := s.db.QueryRow(`
		UPDATE shorturl SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND user_id = $2 RETURNING id, url, code, created_at
	`, id, userId).Scan(&url.ID, &url.URL, &url.Code, &url.CreatedAt)

	if err != nil {
		fmt.Println("Error deleting URL:", err)
		return nil, err
	}

	return &url, nil
}
