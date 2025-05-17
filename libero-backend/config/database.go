package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
	"libero-backend/internal/models"
)

// InitDB initializes the database connection
func InitDB(config *Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	log.Println("Connected to database successfully")
	
	// Auto migrate the database schemas
	migrateDB(db)
	
	return db
}

// migrateDB automatically migrates the database schema
func migrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.CachedFixtures{},
		&models.CachedTodayFixtures{},
		// Add more models here as needed
	)
	
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	
	log.Println("Database migration completed successfully")
}