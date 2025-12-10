package storage

import (
	"context"
	"fmt"
	"rag-training-auth-service/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

type dbConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func Init() (*PostgresDB, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("не собрать конфигурацию подключения к PostgreSQL: %w", err)
	}

	connStr := getConnectString(config)

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("проверка подключения не удалась: %w", err)
	}

	return &PostgresDB{Pool: dbpool}, nil
}

func (db *PostgresDB) Close() error {
	if db.Pool != nil {
		db.Pool.Close()
	}
	return nil
}

func getConnectString(config *dbConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.user, config.password, config.host, config.port, config.dbname,
	)
}

func getConfig() (*dbConfig, error) {
	var config dbConfig
	var err error = nil

	config.host, err = utils.GetEnv("DB_HOST")
	if err != nil {
		return &dbConfig{}, err
	}

	config.port, err = utils.GetEnv("DB_PORT")
	if err != nil {
		return &dbConfig{}, err
	}

	config.user, err = utils.GetEnv("DB_USER")
	if err != nil {
		return &dbConfig{}, err
	}

	config.password, err = utils.GetEnv("DB_PASSWORD")
	if err != nil {
		return &dbConfig{}, err
	}

	config.dbname, err = utils.GetEnv("DB_NAME")
	if err != nil {
		return &dbConfig{}, err
	}

	return &config, err
}
