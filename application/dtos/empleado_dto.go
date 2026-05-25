package dtos

type CreateEmpleadoRequest struct {
	Nombre     string  `json:"nombre" binding:"required"`
	Apellido   string  `json:"apellido" binding:"required"`
	Correo     string  `json:"correo" binding:"required,email"`
	Cargo      string  `json:"cargo" binding:"required"`
	Salario    float64 `json:"salario" binding:"required,gt=0"`
	CompaniaID uint    `json:"compania_id" binding:"required"`
}

type UpdateEmpleadoRequest struct {
	Nombre   string  `json:"nombre" binding:"required"`
	Apellido string  `json:"apellido" binding:"required"`
	Correo   string  `json:"correo" binding:"required,email"`
	Cargo    string  `json:"cargo" binding:"required"`
	Salario  float64 `json:"salario" binding:"required,gt=0"`
}

type EmpleadoResponse struct {
	ID         uint    `json:"id"`
	Nombre     string  `json:"nombre"`
	Apellido   string  `json:"apellido"`
	Correo     string  `json:"correo"`
	Cargo      string  `json:"cargo"`
	Salario    float64 `json:"salario"`
	CompaniaID uint    `json:"compania_id"`
}
