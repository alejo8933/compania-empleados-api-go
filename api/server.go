package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"compania-api/api/handlers"
	"compania-api/api/middlewares"
	"compania-api/api/routes"
	"compania-api/application/services"
	"compania-api/infrastructure/database"
	"compania-api/infrastructure/unit_of_work"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{logger: logger}
}

func (s *Server) Start() {
	s.logger.Info("Iniciando servidor de la API")

	db, err := database.ConnectDB(s.logger)
	if err != nil {
		s.logger.Fatal("No se pudo establecer la conexión a la base de datos", zap.Error(err))
	}

	err = database.SeedDatabase(db, s.logger)
	if err != nil {
		s.logger.Error("Advertencia: Falló la ejecución de semillas", zap.Error(err))
	}

	uow := unit_of_work.NewUnitOfWork(db)
	compService := services.NewCompaniaService(uow, s.logger)
	empService := services.NewEmpleadoService(uow, s.logger)

	compHandler := handlers.NewCompaniaHandler(compService)
	empHandler := handlers.NewEmpleadoHandler(empService)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(middlewares.Logger(s.logger))
	r.Use(middlewares.ErrorHandler(s.logger))

	routes.SetupRoutes(r, compHandler, empHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		s.logger.Info("Servidor HTTP escuchando en el puerto", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Error en ListenAndServe", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.logger.Info("Apagando el servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Fatal("Apagado forzado del servidor:", zap.Error(err))
	}

	s.logger.Info("Servidor apagado de manera exitosa")
}
