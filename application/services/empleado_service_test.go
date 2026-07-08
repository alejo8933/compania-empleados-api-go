package services_test

import (
	"context"
	"testing"
	"compania-api/application/dtos"
	"compania-api/services"
	// Importa tus mocks o usa la base de datos de pruebas configurada
)

func TestRegisterEmpleadoBatch_ShouldRollbackOnError(t *testing.T) {
	// Configura aquí tu entorno de prueba o mock de base de datos
	ctx := context.Background()

	// Creamos un lote donde el segundo registro generará un error intencional (ej: Correo vacío o duplicado)
	lotePrueba := []dtos.EmpleadoBulkInputDTO{
		{Nombre: "EmpleadoValido1", Apellido: "Test", Correo: "valido1@test.com", Salario: 1500000, CompaniaID: 1},
		{Nombre: "EmpleadoInvalido", Apellido: "Test", Correo: "", Salario: 2000000, CompaniaID: 1}, // Correo vacío viola validación
	}

	err := empleadoService.RegisterEmpleadoBatch(ctx, lotePrueba)

	// Verificación 1: El servicio debió retornar un error
	if err == nil {
		t.Errorf("Se esperaba un error en el procesamiento del lote, pero la operación fue exitosa")
	}

	// Verificación 2: Consultar la base de datos para asegurar que "EmpleadoValido1" NO se guardó (Garantía de Rollback)
	empleado, _ := empleadoRepository.GetByEmail("valido1@test.com")
	if empleado != nil {
		t.Errorf("Fallo de atomicidad: El Unit of Work no realizó Rollback. El empleado válido fue guardado inapropiadamente.")
	}
}