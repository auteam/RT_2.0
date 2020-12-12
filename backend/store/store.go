package store

import (
	"database/sql"
	"log"

	//"net/http"
	//"../controller"

	"github.com/BurntSushi/toml"
	_ "github.com/lib/pq" // ...
)

// Config ...
type Config struct {
	BindAddr    string
	DatabaseURL string
	SecretKey   string
}

// Store ...
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	Config         *Config
}

// Configurate ...
func Configurate() *Config {
	var conf Config
	if _, err := toml.DecodeFile("./config/config.toml", &conf); err != nil {
		log.Fatal(err)
	}
	return &conf

}

// NewDB ...
func NewDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
