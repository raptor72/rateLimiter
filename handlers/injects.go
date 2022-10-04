package handlers

import (
	"github.com/jmoiron/sqlx"
	"github.com/raptor72/rateLimiter/api/white_lists"
)

func injectWhiteLists(db *sqlx.DB) *white_lists.Handler {
	pgsqlStorage := white_lists.NewPgsqlStorage(db)
	handler := white_lists.NewHandler(pgsqlStorage)
	return handler
}
