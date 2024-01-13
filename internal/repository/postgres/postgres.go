package postgres

import (
	"github.com/MamushevArup/typeracer/pkg/logger"
	"github.com/MamushevArup/typeracer/pkg/psql"
)

type Repo struct {
	dbconn psql.DBConnector
	lg     *logger.Logger
}

func NewRepository(lg *logger.Logger, dbconn psql.DBConnector) {
	// interface
}
func (r *Repo) Create() {
	_ = r.dbconn.DBConnector()

}
