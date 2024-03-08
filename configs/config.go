package configs

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	User         string
	DbName       string
	Password     string
	Host         string
	Port         int
	Sslmode      string
	MaxOpenConns int
	Timer        uint32
	Db           string
	ServerPort   string
}

var (
	DEFAULT_USER        = "postgres"
	DEFAULT_DBNAME      = "database"
	DEFAULT_PASSWORD    = "password"
	DEFAULT_HOST        = "127.0.0.1"
	DEFAULT_PORT        = "5432"
	DEFAULT_SSLMODE     = "disable"
	DEFAULT_MAXCONNS    = "10"
	DEFAULT_TIMER       = "1"
	DEFAULT_DB          = "postgres"
	DEFAULT_SERVER_PORT = ":8080"
)

func ReadConfig() (*Config, error) {
	user, exists := os.LookupEnv("USER")
	if !exists {
		user = DEFAULT_USER
	}
	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		dbName = DEFAULT_DBNAME
	}
	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		password = DEFAULT_PASSWORD
	}
	host, exists := os.LookupEnv("HOST")
	if !exists {
		host = DEFAULT_HOST
	}
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = DEFAULT_PORT
	}
	port_int, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf(".env port error: %w", err)
	}

	sslmode, exists := os.LookupEnv("SSLMODE")
	if !exists {
		sslmode = DEFAULT_SSLMODE
	}
	maxConns, exists := os.LookupEnv("MAX_OPEN_CONNS")
	if !exists {
		maxConns = DEFAULT_MAXCONNS
	}
	conns_int, err := strconv.Atoi(maxConns)
	if err != nil {
		return nil, fmt.Errorf(".env max conns error: %w", err)
	}

	timer, exists := os.LookupEnv("TIMER")
	if !exists {
		timer = DEFAULT_TIMER
	}
	timer_int, err := strconv.Atoi(timer)
	if err != nil {
		return nil, fmt.Errorf(".env timer error: %w", err)
	}
	db, exists := os.LookupEnv("DB")
	if !exists {
		db = DEFAULT_DB
	}
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		serverPort = DEFAULT_SERVER_PORT
	}

	dbConfig := Config{
		User:         user,
		DbName:       dbName,
		Password:     password,
		Host:         host,
		Port:         port_int,
		Sslmode:      sslmode,
		MaxOpenConns: conns_int,
		Timer:        uint32(timer_int),
		Db:           db,
		ServerPort:   serverPort,
	}

	return &dbConfig, nil
}
