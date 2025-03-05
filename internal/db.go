package internal

import (
	"context"
	"fmt"

	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"go.uber.org/fx"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newDatabaseConnection(lc fx.Lifecycle, config *config.Config) *gorm.DB {
	fmt.Println("NewDatabaseConnection")
	host := config.Postgres.Host
	port := config.Postgres.Port
	user := config.Postgres.User
	password := config.Postgres.Password
	dbName := config.Postgres.DbName

	connectionStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s  sslmode=disable", host, port, user, password, dbName)
	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!\n")
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Connecting to database...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sql, err := db.DB()
			if err != nil {
				err := fmt.Errorf("error while closing database connection: %v", err)
				fmt.Println(err)
			}

			return sql.Close()
		},
	})

	fmt.Println("Connected to database successfully!")
	return db
}
