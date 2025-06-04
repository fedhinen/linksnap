CREATE TABLE shorturl (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    code VARCHAR(10) NOT NULL UNIQUE,
    user_id TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE analytics (
    id SERIAL PRIMARY KEY,
    short_url_id INTEGER NOT NULL,
    ip_address TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (short_url_id) REFERENCES shorturl (id)
);
