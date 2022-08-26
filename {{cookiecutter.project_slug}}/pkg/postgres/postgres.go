package postgres

import (
	"fmt"
	"time"

	"{{cookiecutter.module_path}}/pkg/config"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

// Return new Postgresql db instance.
func NewPsqlDB(c *config.Config) (*sqlx.DB, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.DBName,
		c.Postgres.Password,
	)

	db, err := sqlx.Connect(c.Postgres.PGDriver, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.Postgres.MaxOpenConns)
	db.SetConnMaxLifetime(c.Postgres.ConnMaxLifetime * time.Second)
	db.SetMaxIdleConns(c.Postgres.MaxIdleConns)
	db.SetConnMaxIdleTime(c.Postgres.ConnMaxIdleTime * time.Second)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
