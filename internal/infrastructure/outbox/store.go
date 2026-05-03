package outbox

import (
	"database/sql"

	"github.com/Ali127Dev/xoutbox"
	"github.com/Ali127Dev/xoutbox/postgres"
)

func NewStore(db *sql.DB) (xoutbox.Store[string], error) {
	return postgres.New[string](db), nil
}
