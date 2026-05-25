package database

import (
	"fmt"
	"os"

	"compania-api/domain/entities"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(logger *zap.Logger) (*gorm.DB, error) {
	logger.Info("Cargando variables de entorno desde .env")
	if err := godotenv.Load(); err != nil {
		logger.Warn("No se pudo cargar el archivo .env, usando variables del sistema")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	logger.Info("Iniciando conexión a base de datos PostgreSQL", zap.String("host", dbHost), zap.String("port", dbPort), zap.String("db", dbName))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Error al conectar a la base de datos", zap.Error(err))
		return nil, err
	}

	logger.Info("Conexión exitosa a la base de datos")

	logger.Info("Ejecutando AutoMigrate para entidades Compañía y Empleado")
	err = db.AutoMigrate(&entities.Compania{}, &entities.Empleado{})
	if err != nil {
		logger.Error("Error ejecutando AutoMigrate", zap.Error(err))
		return nil, err
	}

	logger.Info("AutoMigrate ejecutado con éxito")
	return db, nil
}
