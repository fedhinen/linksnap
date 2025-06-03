package models

import "time"

type NewShortUrl struct {
	UserId string
	Url    string
	Code   string
}

type ShortUrl struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id,omitempty"`
}
