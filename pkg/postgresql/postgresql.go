package postgresql

import (
	"database/sql"
	"log"
	"time"

	"github.com/idoyudha/eshop-warehouse/config"
	_ "github.com/lib/pq"
)

const (
	_defaultDriver       = "postgres"
	_defaultConnTimeout  = 2 * time.Second
	_defaultConnAttempts = 4 // (CPU cores × 2)
	_defaultMaxPoolSize  = 10
)

type Postgres struct {
	Conn *sql.DB
}

func NewPostgres(cfg config.PostgreSQL) (*Postgres, error) {
	client, err := sql.Open(_defaultDriver, cfg.URL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = client.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Postgres{Conn: client}, nil
}
