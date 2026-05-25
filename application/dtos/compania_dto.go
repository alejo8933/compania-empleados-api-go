package dtos

import "time"

type CreateCompaniaRequest struct {
	Nombre    string `json:"nombre" binding:"required"`
	Direccion string `json:"direccion" binding:"required"`
	Telefono  string `json:"telefono" binding:"required"`
}

type UpdateCompaniaRequest struct {
	Nombre    string `json:"nombre" binding:"required"`
	Direccion string `json:"direccion" binding:"required"`
	Telefono  string `json:"telefono" binding:"required"`
}

type CompaniaResponse struct {
	ID            uint               `json:"id"`
	Nombre        string             `json:"nombre"`
	Direccion     string             `json:"direccion"`
	Telefono      string             `json:"telefono"`
	FechaCreacion time.Time          `json:"fecha_creacion"`
	Empleados     []EmpleadoResponse `json:"empleados,omitempty"`
}

type CreateCompaniaEmpleadoRequest struct {
	Nombre   string  `json:"nombre" binding:"required"`
	Apellido string  `json:"apellido" binding:"required"`
	Correo   string  `json:"correo" binding:"required,email"`
	Cargo    string  `json:"cargo" binding:"required"`
	Salario  float64 `json:"salario" binding:"required,gt=0"`
}

type CreateCompaniaWithEmpleadosRequest struct {
	Nombre    string                          `json:"nombre" binding:"required"`
	Direccion string                          `json:"direccion" binding:"required"`
	Telefono  string                          `json:"telefono" binding:"required"`
	Empleados []CreateCompaniaEmpleadoRequest `json:"empleados" binding:"required,dive"`
}
