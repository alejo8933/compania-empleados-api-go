package routes

import (
	"compania-api/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, compHandler *handlers.CompaniaHandler, empHandler *handlers.EmpleadoHandler) {
	api := r.Group("/api")
	{
		companias := api.Group("/companias")
		{
			companias.GET("", compHandler.GetAll)
			companias.GET("/:id", compHandler.GetById)
			companias.POST("", compHandler.Create)
			companias.PUT("/:id", compHandler.Update)
			companias.DELETE("/:id", compHandler.Delete)
			companias.GET("/:id/empleados", compHandler.GetEmpleadosByCompania)
			companias.POST("/con-empleados", compHandler.CreateWithEmpleados)
		}

		empleados := api.Group("/empleados")
		{
			empleados.POST("/lote", empHandler.RegisterBatchHandler)
			empleados.GET("/paginado", empHandler.GetPagedHandler)
			empleados.GET("", empHandler.GetAll)
			empleados.GET("/:id", empHandler.GetById)
			empleados.POST("", empHandler.Create)
			empleados.PUT("/:id", empHandler.Update)
			empleados.DELETE("/:id", empHandler.Delete)
			
		}
	}
}
