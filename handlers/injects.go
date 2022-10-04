package handlers

import (
	"github.com/jmoiron/sqlx"
	"github.com/raptor72/rateLimiter/api/white_lists"
)

func injectWhiteLists(db *sqlx.DB) *whitelists.Handler {
	pgsqlStorage := whitelists.NewPgsqlStorage(db)
	handler := whitelists.NewHandler(pgsqlStorage)
	return handler
}
