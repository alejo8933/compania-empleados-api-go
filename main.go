package main

import (
	"log"

	"compania-api/api"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("No se pudo iniciar el logger zap: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Iniciando aplicación de gestión de Compañías y Empleados")

	server := api.NewServer(logger)
	server.Start()
}
