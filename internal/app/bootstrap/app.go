package bootstrap

import (
	"crypto/rand"
	"database/sql"
	"erp/internal/app/models"
	"erp/internal/app/repositories"
	"erp/internal/app/services"
	"erp/internal/infrastructure"
	"erp/internal/pkg/telphin"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	DB          *sql.DB
	Redis       infrastructure.RedisClient
	TelegramBot *tgbotapi.BotAPI
	HTTPServer  *gin.Engine
	Telphin     *telphin.TelphinClient
	Handlers    Handlers
}

type SQLDBAdapter struct {
	db *gorm.DB
}

func NewSQLDBAdapter(gormDB *gorm.DB) *sql.DB {
	sqlDB, _ := gormDB.DB()
	return sqlDB
}

func NewApp() *App {
	gormDB := SetupDatabase()

	// Для старых репозиториев
	sqlDB := NewSQLDBAdapter(gormDB)

	redisClient, err := infrastructure.NewRedisClient()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Fatal("Failed to init Telegram bot:", err)
	}
	telphin, err := telphin.NewTelphinClient(os.Getenv("TELPHIN_API_KEY"))
	if err != nil {
		log.Fatal("Failed to init Telphin client:", err)
	}

	// Репозитории
	roleRepo := repositories.NewRoleRepository(sqlDB)
	userRepo := repositories.NewUserRepository(sqlDB)
	engineerRepo := repositories.NewEngineerRepository(sqlDB)
	orderRepo := repositories.NewOrderRepository(sqlDB)
	callRepo := repositories.NewCallRepository(sqlDB)
	dictRepo := repositories.NewDictionaryRepository(sqlDB)
	reportRepo := repositories.NewReportRepository(sqlDB)

	motivationRepo := repositories.NewMotivationRepository(gormDB)
	engineerTargetRepo := repositories.NewEngineerTargetRepository(gormDB)
	engineerMotivationRepo := repositories.NewEngineerMotivationRepository(gormDB)
	// Сервисы
	blacklist := services.NewTokenBlacklist(redisClient)
	authService := services.NewAuthService(userRepo, blacklist)
	userService := services.NewUserService(userRepo)
	engineerService := services.NewEngineerService(engineerRepo)

	telegramService := services.NewTelegramService(bot)
	callService := services.NewCallService(telphin, callRepo)
	notificationService := services.NewNotificationService(telegramService, callService, redisClient)

	orderService := services.NewOrderService(orderRepo, notificationService)
	dictService := services.NewDictionaryService(dictRepo)
	reportService := services.NewReportService(reportRepo, orderRepo, gormDB)

	motivationService := services.NewMotivationService(motivationRepo)
	egineerTargetService := services.NewEngineerTargetService(engineerTargetRepo)
	engineerMotivationService := services.NewEngineerMotivationService(engineerMotivationRepo)

	if err := ensureAdminExists(userRepo, roleRepo); err != nil {
		log.Printf("WARNING: Admin creation failed: %v", err)
	}

	// Хендлеры

	handlers := NewHandlers(bot, userService, engineerService, orderService, authService, dictService, reportService, motivationService, egineerTargetService, engineerMotivationService)
	httpServer := SetupRouter(handlers)

	return &App{
		DB:          sqlDB,
		Redis:       redisClient,
		TelegramBot: bot,
		HTTPServer:  httpServer,
		Handlers:    handlers,
	}
}

func ensureAdminExists(repo *repositories.UserRepository, roleRepo *repositories.RoleRepository) error {
	username := os.Getenv("DEFAULT_ADMIN_USERNAME")
	password := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	name := os.Getenv("DEFAULT_ADMIN_NAME")

	// Валидация параметров
	if username == "" {
		return fmt.Errorf("DEFAULT_ADMIN_USERNAME not set")
	}

	if password == "" {
		password = generateStrongPassword(16)
		log.Printf("WARNING: Generated admin password: %s", password)
	} else if !isPasswordStrong(password) {
		return fmt.Errorf("admin password does not meet security requirements: %s", password)
	}

	// Проверка существования
	existing, err := repo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			existing = nil
		} else {
			return fmt.Errorf("database error: %w", err)
		}
	}
	if existing != nil {
		return nil
	}

	roleID, err := roleRepo.GetRoleIDByName("admin")
	if err != nil {
		return fmt.Errorf("failed to find 'admin' role: %w", err)
	}

	// Создание
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return fmt.Errorf("password hashing failed: %w", err)
	}

	admin := &models.User{
		Login:     username,
		Password:  string(hashedPass),
		FirstName: name,
		RoleID:    roleID,
	}

	if err := repo.Create(admin); err != nil {
		return fmt.Errorf("create admin failed: %w", err)
	}

	log.Printf("Admin account initialized. Username: %s", username)
	return nil
}

func generateStrongPassword(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" +
		"!@#$%^&*()_+"
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	for i := 0; i < length; i++ {
		buf[i] = chars[int(buf[i])%len(chars)]
	}
	return string(buf)
}

func isPasswordStrong(pwd string) bool {
	return len(pwd) >= 12 &&
		strings.ContainsAny(pwd, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") &&
		strings.ContainsAny(pwd, "abcdefghijklmnopqrstuvwxyz") &&
		strings.ContainsAny(pwd, "0123456789") &&
		strings.ContainsAny(pwd, "!@#$%^&*()_+")
}
