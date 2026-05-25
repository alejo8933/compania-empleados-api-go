// @title			Compania API
// @version		1.0
// @description	API de gestión de compañías y empleados
// @host			localhost:8081
// @BasePath		/
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
