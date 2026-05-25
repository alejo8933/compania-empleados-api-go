package handlers

import (
	"net/http"
	"strconv"

	"compania-api/application/dtos"
	"compania-api/application/services"

	"github.com/gin-gonic/gin"
)

type CompaniaHandler struct {
	service *services.CompaniaService
}

func NewCompaniaHandler(service *services.CompaniaService) *CompaniaHandler {
	return &CompaniaHandler{service: service}
}

// GetAll godoc
//
//	@Summary	Obtener todas las compañías
//	@Tags		compañías
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/api/companias [get]
func (h *CompaniaHandler) GetAll(c *gin.Context) {
	companias, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener las compañías",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    companias,
	})
}

// GetById godoc
//
//	@Summary	Obtener una compañía por ID
//	@Tags		compañías
//	@Produce	json
//	@Param		id	path		int	true	"ID de la compañía"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/api/companias/{id} [get]
func (h *CompaniaHandler) GetById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de compañía inválido",
			"error":   err.Error(),
		})
		return
	}

	compania, err := h.service.GetById(uint(id))
	if err != nil {
		if err.Error() == "compañía no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Compañía no encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener la compañía",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    compania,
	})
}

// Create godoc
//
//	@Summary	Crear una compañía
//	@Tags		compañías
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dtos.CreateCompaniaRequest	true	"Datos de la compañía"
//	@Success	201		{object}	map[string]interface{}
//	@Failure	400		{object}	map[string]interface{}
//	@Failure	500		{object}	map[string]interface{}
//	@Router		/api/companias [post]
func (h *CompaniaHandler) Create(c *gin.Context) {
	var req dtos.CreateCompaniaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	compania, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al crear la compañía",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    compania,
	})
}

// Update godoc
//
//	@Summary	Actualizar una compañía
//	@Tags		compañías
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int							true	"ID de la compañía"
//	@Param		request	body		dtos.UpdateCompaniaRequest	true	"Datos de la compañía"
//	@Success	200		{object}	map[string]interface{}
//	@Failure	400		{object}	map[string]interface{}
//	@Failure	404		{object}	map[string]interface{}
//	@Failure	500		{object}	map[string]interface{}
//	@Router		/api/companias/{id} [put]
func (h *CompaniaHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de compañía inválido",
			"error":   err.Error(),
		})
		return
	}

	var req dtos.UpdateCompaniaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	compania, err := h.service.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "compañía no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Compañía no encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al actualizar la compañía",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    compania,
	})
}

// Delete godoc
//
//	@Summary	Eliminar una compañía
//	@Tags		compañías
//	@Produce	json
//	@Param		id	path	int	true	"ID de la compañía"
//	@Success	204
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/api/companias/{id} [delete]
func (h *CompaniaHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de compañía inválido",
			"error":   err.Error(),
		})
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if err.Error() == "compañía no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Compañía no encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al eliminar la compañía",
			"error":   err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetEmpleadosByCompania godoc
//
//	@Summary	Obtener empleados de una compañía
//	@Tags		compañías
//	@Produce	json
//	@Param		id	path		int	true	"ID de la compañía"
//	@Success	200	{object}	map[string]interface{}
//	@Failure	400	{object}	map[string]interface{}
//	@Failure	404	{object}	map[string]interface{}
//	@Failure	500	{object}	map[string]interface{}
//	@Router		/api/companias/{id}/empleados [get]
func (h *CompaniaHandler) GetEmpleadosByCompania(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID de compañía inválido",
			"error":   err.Error(),
		})
		return
	}

	compania, err := h.service.GetById(uint(id))
	if err != nil {
		if err.Error() == "compañía no encontrada" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Compañía no encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Error al obtener la compañía",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    compania.Empleados,
	})
}

// CreateWithEmpleados godoc
//
//	@Summary	Crear compañía con empleados
//	@Tags		compañías
//	@Accept		json
//	@Produce	json
//	@Param		request	body		dtos.CreateCompaniaWithEmpleadosRequest	true	"Compañía con empleados"
//	@Success	201		{object}	map[string]interface{}
//	@Failure	400		{object}	map[string]interface{}
//	@Router		/api/companias/con-empleados [post]
func (h *CompaniaHandler) CreateWithEmpleados(c *gin.Context) {
	var req dtos.CreateCompaniaWithEmpleadosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Datos de entrada inválidos",
			"error":   err.Error(),
		})
		return
	}

	compania, err := h.service.CreateWithEmpleados(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Error al procesar la creación atómica de compañía y empleados",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Operación exitosa",
		"data":    compania,
	})
}
