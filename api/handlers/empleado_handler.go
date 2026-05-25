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
