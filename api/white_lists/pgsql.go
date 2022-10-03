package white_lists

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const tableName = "white_lists"

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

func (s *PgsqlStorage) Select() ([]*WhiteListModel, error) {
	fields := make([]string, 0, len(tableFields)+1)
	fields = append(fields, "COUNT(id) OVER() as total")
	fields = append(fields, tableFields...)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	selectBuilder := psql.Select(fields...).
		From(tableName).
		OrderBy("id DESC")

	q, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql err: %v", err)
	}

	whiteLists := []*WhiteListModel{}

	if err := s.db.Select(&whiteLists, q, args...); err != nil {
		return nil, fmt.Errorf("Select err: %v", err)
	}

	return whiteLists, nil
}

// func (s *PgsqlStorage) SelectByAddress(address string) ([]*WhiteListModel, error) {
// 	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
// 	selectBuilder := psql.Select(tableFields...).
// 		From(tableName).
// 		Where(sq.Eq{"address": address})

// 	q, args, err := selectBuilder.ToSql()
// 	if err != nil {
// 		return nil, fmt.Errorf("ToSql err: %v", err)
// 	}

// 	whiteLists := []*WhiteListModel{}

// 	if err := s.db.Select(&whiteLists, q, args...); err != nil {
// 		return nil, fmt.Errorf("Select err: %v", err)
// 	}

// 	return whiteLists, nil
// }
