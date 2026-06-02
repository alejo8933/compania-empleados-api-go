package repositories

import (
	"compania-api/domain/entities"
	"context"
	"gorm.io/gorm"
)

type EmpleadoRepositoryImpl struct {
	db *gorm.DB
}

func NewEmpleadoRepository(db *gorm.DB) *EmpleadoRepositoryImpl {
	return &EmpleadoRepositoryImpl{db: db}
}

func (r *EmpleadoRepositoryImpl) GetAll() ([]entities.Empleado, error) {
	var empleados []entities.Empleado
	err := r.db.Find(&empleados).Error
	return empleados, err
}

func (r *EmpleadoRepositoryImpl) GetById(id uint) (*entities.Empleado, error) {
	var e entities.Empleado
	err := r.db.First(&e, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *EmpleadoRepositoryImpl) Create(e *entities.Empleado) error {
	return r.db.Create(e).Error
}

func (r *EmpleadoRepositoryImpl) Update(e *entities.Empleado) error {
	return r.db.Save(e).Error
}

func (r *EmpleadoRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.Empleado{}, id).Error
}

func (r *EmpleadoRepositoryImpl) FindByCondition(query interface{}, args ...interface{}) ([]entities.Empleado, error) {
	var empleados []entities.Empleado
	err := r.db.Where(query, args...).Find(&empleados).Error
	return empleados, err
}
// CreateRange realiza una inserción masiva (Bulk Insert) de empleados en lote.
func (r *EmpleadoRepositoryImpl) CreateRange(ctx context.Context, empleados []*entities.Empleado) error {
	// GORM maneja el lote automáticamente al pasarle el puntero del slice
	if err := r.db.WithContext(ctx).Create(&empleados).Error; err != nil {
		return err
	}
	return nil
}

// GetPaged ejecuta la búsqueda en la BD con paginación, filtros y ordenamiento dinámico.
func (r *EmpleadoRepositoryImpl) GetPaged(ctx context.Context, page int, size int, orderBy string, direction string, search string) ([]entities.Empleado, int64, error) {
	var empleados []entities.Empleado
	var totalRecords int64

	// 1. Inicializar la consulta base sobre la tabla de Empleados
	query := r.db.WithContext(ctx).Model(&entities.Empleado{})

	// 2. Aplicar el Filtro de búsqueda si el cliente envía texto
	if search != "" {
		searchParam := "%" + search + "%"
		query = query.Where("nombre LIKE ? OR apellido LIKE ? OR email LIKE ?", searchParam, searchParam, searchParam)
	}

	// 3. Contar el total de registros en la BD que cumplen con ese filtro
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// 4. Aplicar el Ordenamiento Dinámico de forma segura
	if orderBy != "" {
		if direction != "desc" {
			direction = "asc"
		}
		query = query.Order(orderBy + " " + direction)
	}

	// 5. Calcular la Paginación (límites y saltos de registros)
	if page < 1 { page = 1 }
	if size < 1 { size = 10 }
	offset := (page - 1) * size

	// 6. Ejecutar la consulta final trayendo solo el bloque de datos solicitado
	if err := query.Limit(size).Offset(offset).Find(&empleados).Error; err != nil {
		return nil, 0, err
	}

	return empleados, totalRecords, nil
}