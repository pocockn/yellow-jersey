package database

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm/logger"

	"yellow-jersey/config"
	"yellow-jersey/pkg/logs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	// GormDB holds a database connection.
	GormDB struct {
		maxConnections int
		url            string
	}
)

// NewConnection creates a new connection for the database.
func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	gormDB := GormDB{
		maxConnections: cfg.MaxConnections,
		url:            cfg.DB.ConnectionString(),
	}

	return gormDB.Connect()
}

// Connect connects to the database and passes back the connection so we can
// use it throughout the application
// TODO: Use proper back off or jitter
func (g GormDB) Connect() (*gorm.DB, error) {
	var gormDB *gorm.DB
	var db *sql.DB
	var err error

	for i := 0; i <= 15; i++ {
		gormDB, err = gorm.Open(postgres.Open(g.url), &gorm.Config{})
		if err == nil {
			db, err = gormDB.DB()
			if err != nil {
				return nil, fmt.Errorf("unable to get underlying database %w", err)
			}
			gormDB.Logger.LogMode(logger.Info)
			break
		}

		if i == 15 {
			err = fmt.Errorf("unable to connect to %s after 30 seconds", g.url)
		}

		logs.Logger.Info().Msgf("%d attempt at connecting to the DB", i)
		time.Sleep(2 * time.Second)
	}

	maxConnsPerContainer := g.maxConnections / 4
	db.SetMaxOpenConns(maxConnsPerContainer / 2)
	db.SetConnMaxLifetime(time.Minute * 4)

	logs.Logger.Debug().Msg("successfully connected to database")

	return gormDB, err
}
