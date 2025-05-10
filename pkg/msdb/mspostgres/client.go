package mspostgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "log"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Name     string `mapstructure:"name"`
	Password string `mapstructure:"password"`
	SSL      string `mapstructure:"ssl"`
}

func NewDB(cfg *Config) (*sqlx.DB, func() error, error) {
	dbParams := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Name,
		cfg.SSL,
		cfg.Password,
	)
	log.Println(dbParams)
	db, err := sqlx.Open("postgres", dbParams)
	if err != nil {
		return nil, nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	return db, db.Close, err
}
