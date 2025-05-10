package postgre

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"music-snap/pkg/app"
	"music-snap/pkg/msdb/mspostgres"
	"music-snap/services/musicsnap/internal/config"
	"os"
)

func init() {
	if err := godotenv.Load(MainEnvName); err != nil {
		log.Print(fmt.Sprintf("No '%s' file found", MainEnvName))
	}
}

const MainEnvName = "/Users/abism/Documents/HSE/MusicDay ВКР/musicsnap/.test.env"
const AppCapsName = "MUSICSNAP"
const migPath = "/Users/abism/Documents/HSE/MusicDay ВКР/musicsnap/services/musicsnap/migrations"
const ConfigPath = "/Users/abism/Documents/HSE/MusicDay ВКР/musicsnap/services/musicsnap/config/config.local.yml"

// setupTestDB initializes a test database connection and returns the DB instance, a cleanup function, and a function to clear the database.
func setupTestDB() (*sqlx.DB, func() error, func() error, error) {
	// Load configuration
	configPath := os.Getenv("TCONFIG_" + AppCapsName)
	if configPath == "" {
		fmt.Printf("TCONFIG_%s environment variable is not set\n", AppCapsName)
	}
	cfg, err := config.NewConfig(ConfigPath, "TMUSICSNAP")
	if err != nil {
		return nil, nil, nil, err
	}

	// Initialize the database connection
	db, closeDBHook, err := mspostgres.NewDB(cfg.Postgres)
	if err != nil {
		return nil, nil, nil, app.NewError(0, "initializeDB", "initializeDB", err)
	}

	//run migrations
	err = mspostgres.RunMigrations(migPath, cfg.Postgres, true)
	if err != nil {
		return nil, nil, nil, app.NewError(0, "migrations up", "migrations up", err)
	}

	// Cleanup function to close the database connection
	closeDB := func() error {
		return closeDBHook()
	}

	clearDB := func() error {

		err = mspostgres.RunMigrations(migPath, cfg.Postgres, false)
		if err != nil {
			fmt.Printf(app.NewError(0, "migrations down", "migrations down", err).Error())
		}
		return err

	}

	return db, closeDB, clearDB, nil
}
