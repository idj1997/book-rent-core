package config

import (
	"bufio"
	"github.com/idj1997/book-rent-core/domain"
	golog "log"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenPostgresDB() *gorm.DB {
	DSN := GetPostgresDSN()
	config := GetGormConfig()
	db, err := gorm.Open(postgres.Open(DSN), config)
	if err != nil {
		log.Fatalf("Error while opening DB connection: %v\n", err)
	} else {
		log.Println("Connection opened to DB")
	}

	MigratePostgresDB(db)
	InitPostgresDB(db)
	return db
}

func GetGormConfig() *gorm.Config {
	config := gorm.Config{}
	if ENV == "test" {
		config.Logger = TestLogger()
	}
	return &config
}

func TestLogger() logger.Interface {
	newLogger := logger.New(
		golog.New(os.Stdout, "\r\n", golog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second * 0, // Slow SQL threshold
			LogLevel:      logger.Info,     // Log level
			Colorful:      true,            // Disable color
		})
	return newLogger
}

func ClosePostgresDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting DB from gormDB in closing DB: %v\n", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatalf("Error closing DB: %v\n", err)
	}

	log.Println("Closed DB")
}

func MigratePostgresDB(db *gorm.DB) {
	populateConfig := GetPopulateConfig()
	if populateConfig.Migrate {

		if populateConfig.Init {
			err := db.Migrator().DropTable(&domain.Book{}, &domain.User{}, &domain.RentDetails{})
			if err != nil {
				log.Fatalf("Error while dropping tables: %v", err)
			}
		}

		err := db.AutoMigrate(&domain.Book{}, &domain.User{}, &domain.RentDetails{})
		if err != nil {
			log.Fatalf("Error while migrating DB: %v\n", err)
		} else {
			log.Printf("Migrated DB \n")
		}
	}
}

func InitPostgresDB(db *gorm.DB) {
	populateConfig := GetPopulateConfig()
	if populateConfig.Init {
		statements := LoadStatementsFromFile(populateConfig.File)
		for _, statement := range statements {
			err := db.Exec(statement).Error
			if err != nil {
				log.Fatalf("Error while inserting rows in DB: %v\n", err)
			}
		}
		log.Printf("Populated DB with %v", populateConfig.File)
	}
}

func LoadStatementsFromFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error while loading file %v: %v\n", filename, err)
	}
	defer func() {
		_ = f.Close()
	}()

	s := bufio.NewScanner(f)
	statements := make([]string, 0)
	for s.Scan() {
		statements = append(statements, s.Text())
	}
	return statements
}
