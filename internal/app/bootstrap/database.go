package bootstrap

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func SetupDatabase() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=skip-verify&allowNativePasswords=true&parseTime=true&loc=Europe%%2FMoscow",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Настройте по необходимости
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Получаем *sql.DB для проверки подключения
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sql.DB:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	return db
}
