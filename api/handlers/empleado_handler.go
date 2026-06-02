package handlers

import (
	"net/http"
	"strconv"

	"compania-api/application/dtos"
	"compania-api/application/services"

	"github.com/gin-gonic/gin"
)

type EmpleadoHandler struct {
	service *services.EmpleadoService
}

func NewEmpleadoHandler(service *services.EmpleadoService) *EmpleadoHandler {
	return &EmpleadoHandler{service: service}
}

// GetAll godoc
//
//	@Summary	Obtener todos los empleados
//	@Tags		empleados
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/api/empleados [get]
func (h *EmpleadoHandler) GetAll(c *gin.Context) {
	empleados, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener los empleados",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    empleados,
	})
}

// GetById godoc
//
//	@Summary	Obtener un empleado por ID
//	@Tags		empleados
//	@Produce	json
//	@Param		id	path		int	true	"ID del empleado"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/api/empleados/{id} [get]
func (h *EmpleadoHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de empleado inválido",
			"error":   err.Error(),
		})
		return
	}

	empleado, err := h.service.GetById(uint(id))
	if err != nil {
		if err.Error() == "empleado no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Empleado no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener el empleado",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    empleado,
	})
}

// Create godoc
//
//	@Summary	Crear un empleado
//	@Tags		empleados
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dtos.CreateEmpleadoRequest	true	"Datos del empleado"
//	@Success	201		{object}	map[string]interface{}
//	@Failure	400		{object}	map[string]interface{}
//	@Router		/api/empleados [post]
func (h *EmpleadoHandler) Create(c *gin.Context) {
	var req dtos.CreateEmpleadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	empleado, err := h.service.Create(&req)
	if err != nil {
		if err.Error() == "compañía asociada no encontrada" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error al crear el empleado",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    empleado,
	})
}

// Update godoc
//
//	@Summary	Actualizar un empleado
//	@Tags		empleados
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int							true	"ID del empleado"
//	@Param		request	body		dtos.UpdateEmpleadoRequest	true	"Datos del empleado"
//	@Success	200		{object}	map[string]interface{}
//	@Failure	400		{object}	map[string]interface{}
//	@Failure	404		{object}	map[string]interface{}
//	@Failure	500		{object}	map[string]interface{}
//	@Router		/api/empleados/{id} [put]
func (h *EmpleadoHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de empleado inválido",
			"error":   err.Error(),
		})
		return
	}

	var req dtos.UpdateEmpleadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	empleado, err := h.service.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "empleado no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Empleado no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al actualizar el empleado",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    empleado,
	})
}

// Delete godoc
//
//	@Summary	Eliminar un empleado
//	@Tags		empleados
//	@Produce	json
//	@Param		id	path	int	true	"ID del empleado"
//	@Success	204
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/api/empleados/{id} [delete]
func (h *EmpleadoHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de empleado inválido",
			"error":   err.Error(),
		})
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if err.Error() == "empleado no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Empleado no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al eliminar el empleado",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
// RegisterBatchHandler procesa la inserción masiva aplicando validaciones (POST /api/empleados/lote)
func (h *EmpleadoHandler) RegisterBatchHandler(c *gin.Context) {
	var inputDTOs []dtos.EmpleadoBulkInputDTO

	// 1. Decodificar el JSON usando la herramienta nativa de Gin
	if err := c.ShouldBindJSON(&inputDTOs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}

	// 2. VALIDACIÓN: Enviar a validar si tus DTOs lo requieren, o procesar directo
	// Llamamos al servicio de aplicación que ya dejamos listo y aprobado
	err := h.service.RegisterEmpleadoBatch(c.Request.Context(), inputDTOs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar el lote: " + err.Error()})
		return
	}

	// 3. Responder con estatus 201 Created según exige la guía
	c.JSON(http.StatusCreated, gin.H{"message": "Lote de empleados registrado exitosamente"})
}

// GetPagedHandler procesa el listado avanzado con filtros (GET /api/empleados/paginado)
func (h *EmpleadoHandler) GetPagedHandler(c *gin.Context) {
	// Capturar parámetros de la URL usando los métodos de Gin
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	orderBy := c.Query("orderBy")
	direction := c.Query("direction")
	search := c.Query("search")

	empleados, total, err := h.service.GetEmpleadoPaged(c.Request.Context(), page, size, orderBy, direction, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener empleados: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":          empleados,
		"total_records": total,
		"page":          page,
		"size":          size,
	})
}