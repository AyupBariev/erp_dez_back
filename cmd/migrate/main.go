package main

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv"
)

/*
  - Накатить миграции
    go run cmd/migrate/main.go up
    *
  - откатить одну миграцию
    go run cmd/migrate/main.go down --steps=1
    *
  - Откатить все миграции
    go run cmd/migrate/main.go down
    *
  - Принудительно установить версию
    go run cmd/migrate/main.go force 3
*/
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <command> [--steps=N]")
		fmt.Println("Commands: up | down | force <version>")
		os.Exit(1)
	}

	if os.Getenv("APP_ENV") != "production" {
		_ = godotenv.Load(".env")
	}

	// Загружаем .env, если нужно
	// можно использовать godotenv, если у тебя переменные только в .env
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_DATABASE")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "mysql" // значение по умолчанию для docker-compose
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306" // порт внутри docker сети
	}

	// Формируем URL для mysql
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Путь до миграций
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		log.Fatal("Failed to init migrate:", err)
	}
	defer m.Close()

	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migrations applied successfully.")

	case "down":
		stepsFlag := flag.NewFlagSet("down", flag.ExitOnError)
		steps := stepsFlag.Int("steps", 0, "Number of steps to rollback (0 = all)")
		_ = stepsFlag.Parse(os.Args[2:])

		if *steps > 0 {
			if err := m.Steps(-*steps); err != nil && err != migrate.ErrNoChange {
				log.Fatal(err)
			}
			fmt.Printf("Rolled back %d step(s) successfully.\n", *steps)
		} else {
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatal(err)
			}
			fmt.Println("Rolled back all migrations successfully.")
		}

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Usage: migrate force <version>")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Invalid version number:", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Forced migration version to %d.\n", version)

	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Usage: migrate create <name>")
		}
		name := os.Args[2]
		// каталог миграций
		dir := "migrations"

		// получаем текущий timestamp для имени
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)

		// формируем имена файлов
		upFile := fmt.Sprintf("%s/%s_%s.up.sql", dir, timestamp, name)
		downFile := fmt.Sprintf("%s/%s_%s.down.sql", dir, timestamp, name)

		// создаём пустые файлы
		if err := os.WriteFile(upFile, []byte("-- +migrate Up\n\n"), 0644); err != nil {
			log.Fatal("Failed to create up file:", err)
		}
		if err := os.WriteFile(downFile, []byte("-- +migrate Down\n\n"), 0644); err != nil {
			log.Fatal("Failed to create down file:", err)
		}

		fmt.Printf("Created new migration files:\n%s\n%s\n", upFile, downFile)

	default:
		fmt.Println("Unknown command:", cmd)
		os.Exit(1)
	}
}
