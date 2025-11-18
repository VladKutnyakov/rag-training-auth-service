package postgres

import (
	"context"
	"fmt"

	"rag-training-auth-service/internal/utils"

	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	Conn *pgx.Conn
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
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("проверка подключения не удалась: %w", err)
	}

	return &PostgresDB{Conn: conn}, nil
}

func (db *PostgresDB) Close() error {
	if db.Conn != nil {
		return db.Conn.Close(context.Background())
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
