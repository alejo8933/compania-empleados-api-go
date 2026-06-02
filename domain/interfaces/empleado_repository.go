package interfaces

import (
	"compania-api/domain/entities"
	"context"
)

type IEmpleadoRepository interface {
	GetAll() ([]entities.Empleado, error)
	GetById(id uint) (*entities.Empleado, error)
	Create(e *entities.Empleado) error
	Update(e *entities.Empleado) error
	Delete(id uint) error
	FindByCondition(query interface{}, args ...interface{}) ([]entities.Empleado, error)
	
	// ==========================================================
	// REQUERIMIENTOS MÓDULO 1 - PARTE II
	// ==========================================================
	CreateRange(ctx context.Context, empleados []*entities.Empleado) error
	GetPaged(ctx context.Context, page int, size int, orderBy string, direction string, search string) ([]entities.Empleado, int64, error)
}