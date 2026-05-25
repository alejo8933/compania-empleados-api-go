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

	// 204 No Content
	c.Status(http.StatusNoContent)
}

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
		// Podría ser un correo duplicado u otro error de validación de base de datos
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
