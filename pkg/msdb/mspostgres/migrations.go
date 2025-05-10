package mspostgres

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // THIS IS CRUCIAL
	"music-snap/pkg/app"
)

func RunMigrations(migrationDirPath string, cfg *Config, up bool) error {
	dbURL :=
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			//fmt.Sprintf("postgres://mattes:secret@localhost:5432/database?sslmode=disable",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
			cfg.SSL,
		)
	m, err := migrate.New("file://"+migrationDirPath, dbURL)
	if err != nil {
		return app.NewError(0, fmt.Sprintf("migrate.New() error %s", dbURL),
			fmt.Sprintf("migrate.New() error %s", dbURL), err)
	}
	if up {
		return m.Up()
	} else {
		return m.Down()
	}
}
