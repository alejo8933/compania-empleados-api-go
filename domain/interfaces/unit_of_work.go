package interfaces

type IUnitOfWork interface {
	Companias() ICompaniaRepository
	Empleados() IEmpleadoRepository
	BeginTransaction() error
	Commit() error
	Rollback() error
	Save() error
}
