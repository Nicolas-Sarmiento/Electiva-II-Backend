package repository

import (
	"fmt"
	"log"

	"ancianato-backend/internal/config"
	"ancianato-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB inicializa la conexión con la base de datos PostgreSQL
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	gormConfig := &gorm.Config{}

	if cfg.DBDebug == "true" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	log.Println("Conexión a PostgreSQL establecida exitosamente")

	// Ejecutar las migraciones
	err = runMigrations(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// runMigrations hace que GORM cree o actualice las tablas según nuestras estructuras en domain.
func runMigrations(db *gorm.DB) error {
	log.Println("Ejecutando AutoMigrate...")
	err := db.AutoMigrate(
		&domain.Room{},
		&domain.MedicalCondition{},
		&domain.EmergencyContact{},
		&domain.Patient{},
		&domain.PatientCondition{},
		&domain.PatientContact{},
		&domain.Wearable{},
		&domain.PatientWearable{},
		&domain.Shift{},
		&domain.Nurse{},
		&domain.AlertType{},
		&domain.Alert{},
	)
	if err != nil {
		return fmt.Errorf("error al migrar la base de datos: %w", err)
	}
	log.Println("Migraciones ejecutadas correctamente")
	return nil
}
