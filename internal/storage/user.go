package storage

import (
	"context"
	"strings"

	"github.com/didsqq/news_feed_bot/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type dbUser struct {
	ID       int64  `db:"id"`
	ChatID   int64  `db:"chat_id"`
	Keywords string `db:"keywords"`
}

type UserPostgresStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserPostgresStorage {
	return &UserPostgresStorage{db: db}
}

func (s *UserPostgresStorage) GetAll(ctx context.Context) ([]model.User, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var users []dbUser
	if err := conn.SelectContext(ctx, &users, `SELECT * FROM users`); err != nil {
		return nil, err
	}

	return lo.Map(users, func(user dbUser, _ int) model.User {
		return model.User{
			ChatID:   user.ChatID,
			Keywords: strings.Split(user.Keywords, ";"),
		}
	}), nil
}

func (s *UserPostgresStorage) Add(ctx context.Context, user model.User) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `INSERT INTO users (chat_id, keywords) VALUES ($1, $2)`, user.ChatID, ""); err != nil {
		return err
	}

	return nil
}

func (s *UserPostgresStorage) AddKeywords(ctx context.Context, user model.User) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i, keyword := range user.Keywords {
		user.Keywords[i] = strings.ToLower(keyword)
	}

	if _, err := conn.ExecContext(ctx, `UPDATE users keywords = $1 WHERE chat_id = $2`,
		strings.Join(user.Keywords, ";"), user.ChatID); err != nil {
		return err
	}

	return nil
}

func (s *UserPostgresStorage) GetKeywords(ctx context.Context, chatId int64) (string, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return "", err
	}

	var user dbUser
	if err := conn.GetContext(ctx, &user, `SELECT * FROM users WHERE chat_id = $1`, chatId); err != nil {
		return "", err
	}

	return user.Keywords, nil
}

func (s *UserPostgresStorage) Delete(ctx context.Context, chatId int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}

	if _, err := conn.ExecContext(ctx, `DELETE FROM users WHERE chat_id = $1;`, chatId); err != nil {
		return err
	}

	return nil
}
