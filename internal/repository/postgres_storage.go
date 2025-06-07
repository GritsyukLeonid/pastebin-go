package repository

import (
	"context"
	"database/sql"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

// Paste
func (s *PostgresStorage) SavePaste(p model.Paste) error {
	query := `INSERT INTO pastes (id, hash, content, created_at, expires_at, views) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(context.Background(), query, p.ID, p.Hash, p.Content, p.CreatedAt, p.ExpiresAt, p.Views)
	return err
}

func (s *PostgresStorage) GetPasteByID(id string) (*model.Paste, error) {
	query := `SELECT id, hash, content, created_at, expires_at, views FROM pastes WHERE id = $1`
	row := s.db.QueryRowContext(context.Background(), query, id)
	var p model.Paste
	err := row.Scan(&p.ID, &p.Hash, &p.Content, &p.CreatedAt, &p.ExpiresAt, &p.Views)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PostgresStorage) DeletePaste(id string) error {
	query := `DELETE FROM pastes WHERE id = $1`
	_, err := s.db.ExecContext(context.Background(), query, id)
	return err
}

func (s *PostgresStorage) DeleteExpiredPastes() error {
	query := `DELETE FROM pastes WHERE expires_at < NOW()`
	_, err := s.db.ExecContext(context.Background(), query)
	return err
}

func (s *PostgresStorage) GetAllPastes() ([]model.Paste, error) {
	query := `SELECT id, hash, content, created_at, expires_at, views FROM pastes`
	rows, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pastes []model.Paste
	for rows.Next() {
		var p model.Paste
		err := rows.Scan(&p.ID, &p.Hash, &p.Content, &p.CreatedAt, &p.ExpiresAt, &p.Views)
		if err != nil {
			return nil, err
		}
		pastes = append(pastes, p)
	}
	return pastes, nil
}

// User
func (s *PostgresStorage) SaveUser(u model.User) error {
	query := `INSERT INTO users (id, username) VALUES ($1, $2)`
	_, err := s.db.ExecContext(context.Background(), query, u.ID, u.Username)
	return err
}

func (s *PostgresStorage) GetUserByID(id string) (*model.User, error) {
	query := `SELECT id, username FROM users WHERE id = $1`
	row := s.db.QueryRowContext(context.Background(), query, id)
	var u model.User
	err := row.Scan(&u.ID, &u.Username)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *PostgresStorage) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.ExecContext(context.Background(), query, id)
	return err
}

func (s *PostgresStorage) GetAllUsers() ([]model.User, error) {
	query := `SELECT id, username FROM users`
	rows, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.ID, &u.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// ShortURL
func (s *PostgresStorage) SaveShortURL(u model.ShortURL) error {
	query := `INSERT INTO shorturls (id, original) VALUES ($1, $2)`
	_, err := s.db.ExecContext(context.Background(), query, u.ID, u.Original)
	return err
}

func (s *PostgresStorage) GetShortURLByID(id string) (*model.ShortURL, error) {
	query := `SELECT id, original FROM shorturls WHERE id = $1`
	row := s.db.QueryRowContext(context.Background(), query, id)
	var u model.ShortURL
	err := row.Scan(&u.ID, &u.Original)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *PostgresStorage) DeleteShortURL(id string) error {
	query := `DELETE FROM shorturls WHERE id = $1`
	_, err := s.db.ExecContext(context.Background(), query, id)
	return err
}

func (s *PostgresStorage) GetAllShortURLs() ([]model.ShortURL, error) {
	query := `SELECT id, original FROM shorturls`
	rows, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []model.ShortURL
	for rows.Next() {
		var u model.ShortURL
		err := rows.Scan(&u.ID, &u.Original)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	return urls, nil
}

// Stats
func (s *PostgresStorage) SaveStats(st model.Stats) error {
	query := `INSERT INTO stats (id, views) VALUES ($1, $2)`
	_, err := s.db.ExecContext(context.Background(), query, st.ID, st.Views)
	return err
}

func (s *PostgresStorage) GetStatsByID(id string) (*model.Stats, error) {
	query := `SELECT id, views FROM stats WHERE id = $1`
	row := s.db.QueryRowContext(context.Background(), query, id)
	var st model.Stats
	err := row.Scan(&st.ID, &st.Views)
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (s *PostgresStorage) DeleteStats(id string) error {
	query := `DELETE FROM stats WHERE id = $1`
	_, err := s.db.ExecContext(context.Background(), query, id)
	return err
}

func (s *PostgresStorage) GetAllStats() ([]model.Stats, error) {
	query := `SELECT id, views FROM stats`
	rows, err := s.db.QueryContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []model.Stats
	for rows.Next() {
		var st model.Stats
		err := rows.Scan(&st.ID, &st.Views)
		if err != nil {
			return nil, err
		}
		stats = append(stats, st)
	}
	return stats, nil
}
