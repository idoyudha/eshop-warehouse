package postgresql

import (
	"database/sql"
	"fmt"
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
	connStr := fmt.Sprintf("%s://%s", _defaultDriver, cfg.URL)

	client, err := sql.Open(_defaultDriver, connStr)
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
