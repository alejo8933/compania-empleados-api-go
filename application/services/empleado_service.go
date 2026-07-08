package services

import (
	"context"
	"compania-api/application/dtos"
	"compania-api/domain/entities"
	"compania-api/domain/interfaces"
	"errors"
	
	"go.uber.org/zap"
)

type EmpleadoService struct {
	uow    interfaces.IUnitOfWork
	logger *zap.Logger
}

func NewEmpleadoService(uow interfaces.IUnitOfWork, logger *zap.Logger) *EmpleadoService {
	return &EmpleadoService{
		uow:    uow,
		logger: logger,
	}
}

func (s *EmpleadoService) GetAll() ([]dtos.EmpleadoResponse, error) {
	empleados, err := s.uow.Empleados().GetAll()
	if err != nil {
		s.logger.Error("Error al obtener todos los empleados", zap.Error(err))
		return nil, err
	}

	var response []dtos.EmpleadoResponse
	for _, e := range empleados {
		response = append(response, dtos.EmpleadoResponse{
			ID:         e.ID,
			Nombre:     e.Nombre,
			Apellido:   e.Apellido,
			Correo:     e.Correo,
			Cargo:      e.Cargo,
			Salario:    e.Salario,
			CompaniaID: e.CompaniaID,
		})
	}
	return response, nil
}

func (s *EmpleadoService) GetById(id uint) (*dtos.EmpleadoResponse, error) {
	e, err := s.uow.Empleados().GetById(id)
	if err != nil {
		s.logger.Error("Error al obtener empleado por ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	if e == nil {
		return nil, errors.New("empleado no encontrado")
	}

	return &dtos.EmpleadoResponse{
		ID:         e.ID,
		Nombre:     e.Nombre,
		Apellido:   e.Apellido,
		Correo:     e.Correo,
		Cargo:      e.Cargo,
		Salario:    e.Salario,
		CompaniaID: e.CompaniaID,
	}, nil
}

func (s *EmpleadoService) Create(req *dtos.CreateEmpleadoRequest) (*dtos.EmpleadoResponse, error) {
	// Verificar si la compañía existe
	c, err := s.uow.Companias().GetById(req.CompaniaID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("compañía asociada no encontrada")
	}

	e := &entities.Empleado{
		Nombre:     req.Nombre,
		Apellido:   req.Apellido,
		Correo:     req.Correo,
		Cargo:      req.Cargo,
		Salario:    req.Salario,
		CompaniaID: req.CompaniaID,
	}

	err = s.uow.Empleados().Create(e)
	if err != nil {
		s.logger.Error("Error al crear empleado", zap.String("correo", req.Correo), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Empleado creado",
		zap.String("nombre", e.Nombre),
		zap.String("compañía", c.Nombre),
	)

	return &dtos.EmpleadoResponse{
		ID:         e.ID,
		Nombre:     e.Nombre,
		Apellido:   e.Apellido,
		Correo:     e.Correo,
		Cargo:      e.Cargo,
		Salario:    e.Salario,
		CompaniaID: e.CompaniaID,
	}, nil
}

func (s *EmpleadoService) Update(id uint, req *dtos.UpdateEmpleadoRequest) (*dtos.EmpleadoResponse, error) {
	e, err := s.uow.Empleados().GetById(id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, errors.New("empleado no encontrado")
	}

	e.Nombre = req.Nombre
	e.Apellido = req.Apellido
	e.Correo = req.Correo
	e.Cargo = req.Cargo
	e.Salario = req.Salario

	err = s.uow.Empleados().Update(e)
	if err != nil {
		s.logger.Error("Error al actualizar empleado", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Empleado actualizado", zap.Uint("id", e.ID))

	return &dtos.EmpleadoResponse{
		ID:         e.ID,
		Nombre:     e.Nombre,
		Apellido:   e.Apellido,
		Correo:     e.Correo,
		Cargo:      e.Cargo,
		Salario:    e.Salario,
		CompaniaID: e.CompaniaID,
	}, nil
}

func (s *EmpleadoService) Delete(id uint) error {
	e, err := s.uow.Empleados().GetById(id)
	if err != nil {
		return err
	}
	if e == nil {
		return errors.New("empleado no encontrado")
	}

	err = s.uow.Empleados().Delete(id)
	if err != nil {
		s.logger.Error("Error al eliminar empleado", zap.Uint("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Empleado eliminado", zap.Uint("id", id))
	return nil
}
// RegisterEmpleadoBatch orquesta la inserción masiva de empleados usando transaccionalidad atómica del UoW
func (s *EmpleadoService) RegisterEmpleadoBatch(ctx context.Context, dtos []dtos.EmpleadoBulkInputDTO) error {
	// 1. Iniciar la transacción formal desde el Unit of Work
	tx, err := s.uow.BeginTx(ctx)
	if err != nil {
		s.logger.Error("No se pudo iniciar la transacción del lote", zap.Error(err))
		return err
	}

	// Asegurar el Rollback automático en caso de un pánico imprevisto en el hilo de ejecución
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var entidades []*entities.Empleado
	for _, dto := range dtos {
		empleado := &entities.Empleado{
			Nombre:     dto.Nombre,
			Apellido:   dto.Apellido,
			Correo:     dto.Correo,
			Cargo:      dto.Cargo,
			Salario:    dto.Salario,
			CompaniaID: dto.CompaniaID,
		}
		entidades = append(entidades, empleado)
	}

	// 2. Guardar el lote completo. Si falla un solo registro, se ejecuta Rollback
	if err := s.uow.Empleados().CreateRange(ctx, entidades); err != nil {
		s.logger.Warn("Error detectado en el lote de empleados. Ejecutando Rollback transaccional...", zap.Error(err))
		tx.Rollback() 
		return err
	}

	// 3. Si todos los registros son válidos, se consolidan permanentemente los cambios en PostgreSQL
	if err := tx.Commit(); err != nil {
		return err
	}

	s.logger.Info("Lote de empleados insertado con éxito de forma atómica")
	return nil
}
// GetEmpleadoPaged orquesta la consulta avanzada con filtros y paginación
func (s *EmpleadoService) GetEmpleadoPaged(ctx context.Context, page, size int, orderBy, direction, search string) ([]entities.Empleado, int64, error) {
	return s.uow.Empleados().GetPaged(ctx, page, size, orderBy, direction, search)
}