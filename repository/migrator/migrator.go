package migrator

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func New(client *mongo.Client, dbConfig *mongodb.Config, msource string) (*Migrator, error) {
	driver, err := mongodb.WithInstance(client, dbConfig)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(msource, dbConfig.DatabaseName, driver)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		migrate: m,
	}, nil
}

func (m *Migrator) Up() error {
	return m.migrate.Up()
}

func (m *Migrator) Down() error {
	return m.migrate.Down()
}
