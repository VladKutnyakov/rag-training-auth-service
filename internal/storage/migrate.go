package storage

import (
	"log"
	"rag-training-auth-service/migrations"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigrations(db *PostgresDB) {
	dbStd := stdlib.OpenDBFromPool(db.Pool)

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		log.Fatal(err)
	}

	goose.SetBaseFS(migrations.Embed)

	if err := goose.Up(dbStd, "."); err != nil {
		log.Fatal(err)
	}

	goose.Status(dbStd, ".")

	log.Println("Миграции применены")

	db.Close()
}
