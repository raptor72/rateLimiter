package blacklists

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const tableName = "black_lists"

var tableFields = []string{
	"id",
	"address",
}

type PgsqlStorage struct {
	db *sqlx.DB
}

func NewPgsqlStorage(db *sqlx.DB) *PgsqlStorage {
	return &PgsqlStorage{db: db}
}

func (s *PgsqlStorage) Select() ([]*BlackListModel, error) {
	fields := make([]string, 0, len(tableFields)+1)
	fields = append(fields, "COUNT(id) OVER() as total")
	fields = append(fields, tableFields...)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	selectBuilder := psql.Select(fields...).
		From(tableName).
		OrderBy("id DESC")

	q, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("toSql err: %w", err)
	}

	whiteLists := []*BlackListModel{}

	if err := s.db.Select(&whiteLists, q, args...); err != nil {
		return nil, fmt.Errorf("select err: %w", err)
	}

	return whiteLists, nil
}
