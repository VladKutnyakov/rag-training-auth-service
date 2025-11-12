package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
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
	config := getConfig()
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

func getConfig() *dbConfig {
	_ = godotenv.Load()

	var config dbConfig

	config.host = getEnv("DB_HOST", "localhost")
	config.port = getEnv("DB_PORT", "5432")
	config.user = getEnv("DB_USER", "postgres")
	config.password = getEnv("DB_PASSWORD", "postgres")
	config.dbname = getEnv("DB_NAME", "myapp")

	return &config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
