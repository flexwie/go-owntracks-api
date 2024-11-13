package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

var WithDb = fx.Provide(
	newDb,
)

var schema = `
CREATE TABLE IF NOT EXISTS location (
	id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
	created TIMESTAMP DEFAULT NOW(),
	lat REAL,
	lng REAL,
	alt REAL
);

ALTER TABLE location ADD IF NOT EXISTS vel REAL NOT NULL DEFAULT 0;
`

func newDb(lc fx.Lifecycle) (*sqlx.DB, error) {
	host, err := checkOrError("db-host")
	user, err := checkOrError("db-user")
	pwd, err := checkOrError("db-password")
	dbn, err := checkOrError("db-name")

	if err != nil {
		return nil, err
	}

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, pwd, dbn)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := db.Close(); err != nil {
				return err
			}

			return nil
		},
	})

	return db, nil
}

func checkOrError(key string) (string, error) {
	value := viper.GetString(key)

	if value == "" {
		return "", errors.New(fmt.Sprintf("%s is not set", key))
	}

	return value, nil
}
