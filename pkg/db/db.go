package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rzkhosroshahi/velox/config"
	"github.com/rzkhosroshahi/velox/pkg/logger"
)

func New(conf *config.DataBaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Name, conf.SSLMODE,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("connected to the database!")

	return db, nil
}
