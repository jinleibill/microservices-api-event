package base

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
)

type Client interface {
	Exec(ctx context.Context, sql string, args ...any) error
	Query(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) *sql.Row
	Migrate(tableName string, sql string) error
	GetTx() *gorm.DB
}

var _ Client = (*sessionClient)(nil)

type sessionClient struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewSessionClient(db *gorm.DB) Client {
	return &sessionClient{db: db}
}

func (s *sessionClient) Exec(ctx context.Context, sql string, args ...any) error {
	db := s.tx.WithContext(ctx).Exec(sql, args...)
	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (s *sessionClient) Query(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	return s.tx.WithContext(ctx).Raw(sql, args...).Rows()
}

func (s *sessionClient) QueryRow(ctx context.Context, sql string, args ...any) *sql.Row {
	return s.tx.WithContext(ctx).Raw(sql, args...).Row()
}

func (s *sessionClient) Migrate(tableName string, sql string) error {
	if !s.db.Migrator().HasTable(tableName) {
		return s.db.Exec(fmt.Sprintf(sql, tableName)).Error
	}

	return nil
}

func (s *sessionClient) GetTx() *gorm.DB {
	s.tx = s.db.Begin()

	return s.tx
}
