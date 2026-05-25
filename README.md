# API de Compañías y Empleados

API REST construida en Go aplicando Onion Architecture, Repository Pattern y Unit of Work para la gestión de compañías y empleados.

---

## Tecnología usada

| Herramienta | Descripción |
|---|---|
| Go | Lenguaje de programación principal |
| Gin | Framework web para crear la API REST |
| GORM | ORM para manejo de base de datos |
| PostgreSQL | Motor de base de datos relacional |
| godotenv | Carga de variables de entorno desde `.env` |
| zap | Logging estructurado |

---

## ORM usado

Se usó **GORM** como ORM principal. Es el más popular en el ecosistema Go, permite trabajar con entidades, relaciones, migraciones y transacciones de forma sencilla y se integra bien con PostgreSQL.

---

## Arquitectura aplicada

El proyecto aplica **Onion Architecture**, donde el dominio está en el centro y las dependencias siempre apuntan hacia adentro. Esto protege la lógica de negocio de detalles técnicos como el ORM o el framework web.

```text
Presentation (API)
    └── Application (Services / DTOs)
            └── Domain (Entities / Interfaces)
    Infrastructure (ORM / DB / Repositories)
```

---

## Estructura del proyecto

```text
compania-api/
├── api/
│   ├── handlers/
│   │   ├── compania_handler.go
│   │   └── empleado_handler.go
│   ├── routes/
│   │   └── routes.go
│   └── middlewares/
├── application/
│   ├── services/
│   │   ├── compania_service.go
│   │   └── empleado_service.go
│   └── dtos/
│       ├── compania_dto.go
│       └── empleado_dto.go
├── domain/
│   ├── entities/
│   │   ├── compania.go
│   │   └── empleado.go
│   └── interfaces/
│       ├── compania_repository.go
│       ├── empleado_repository.go
│       └── unit_of_work.go
├── infrastructure/
│   ├── database/
│   │   ├── connection.go
│   │   ├── migrations.go
│   │   └── seed.go
│   ├── repositories/
│   │   ├── compania_repository_impl.go
│   │   └── empleado_repository_impl.go
│   └── unit_of_work/
│       └── unit_of_work_impl.go
├── .env
├── go.mod
├── go.sum
└── main.go
```

---

## Entidades

### Compañía

| Campo | Tipo | Descripción |
|---|---|---|
| Id | uint | Llave primaria |
| Nombre | string | Nombre de la compañía |
| Direccion | string | Dirección física |
| Telefono | string | Número de contacto |
| FechaCreacion | time.Time | Fecha de registro |

### Empleado

| Campo | Tipo | Descripción |
|---|---|---|
| Id | uint | Llave primaria |
| Nombre | string | Nombre del empleado |
| Apellido | string | Apellido del empleado |
| Correo | string | Correo electrónico |
| Cargo | string | Cargo o rol |
| Salario | float64 | Salario asignado |
| CompaniaId | uint | Llave foránea hacia Compañía |

---

## Relación entre entidades

```text
Compañía 1 ──── * Empleado
```

Una compañía puede tener muchos empleados. Cada empleado pertenece a una sola compañía.

---

## Repository Pattern

Se definieron interfaces en la capa `domain` y sus implementaciones concretas en `infrastructure`. Esto permite que la lógica de negocio no dependa directamente del ORM o de la base de datos.

Cada repositorio expone los métodos:

- `GetAll` — Obtener todos los registros.
- `GetById` — Obtener un registro por id.
- `Create` — Crear un nuevo registro.
- `Update` — Actualizar un registro existente.
- `Delete` — Eliminar un registro.

---

## Unit of Work

### ¿Qué es Unit of Work?

Unit of Work es un patrón de diseño que agrupa varias operaciones de base de datos en una sola unidad lógica de trabajo. Coordina repositorios y controla una única sesión de base de datos para garantizar que todos los cambios se confirmen o reviertan juntos.

### ¿Cómo se implementó en esta tecnología?

Se implementó una interfaz `UnitOfWork` en `domain/interfaces` y su implementación concreta en `infrastructure/unit_of_work`. La implementación usa una instancia de `*gorm.DB` para controlar la transacción.

### ¿Cómo se manejan las transacciones?

Las transacciones se manejan a través de `*gorm.DB`. Al iniciar una operación que involucra varios repositorios, se abre una transacción que envuelve todas las operaciones.

### ¿Cómo se hace commit?

Se llama al método `Commit()` de la instancia de transacción cuando todas las operaciones finalizan correctamente.

### ¿Cómo se hace rollback?

Se llama al método `Rollback()` cuando ocurre cualquier error durante la transacción, deshaciendo todos los cambios realizados hasta ese punto.

---

## Endpoints

### Compañías

| Método | Ruta | Descripción |
|---|---|---|
| GET | /api/companias | Listar todas las compañías |
| GET | /api/companias/:id | Consultar una compañía por id |
| POST | /api/companias | Crear una compañía |
| PUT | /api/companias/:id | Actualizar una compañía |
| DELETE | /api/companias/:id | Eliminar una compañía |
| GET | /api/companias/:id/empleados | Listar empleados de una compañía |
| POST | /api/companias/con-empleados | Crear compañía con empleados (transaccional) |

### Empleados

| Método | Ruta | Descripción |
|---|---|---|
| GET | /api/empleados | Listar todos los empleados |
| GET | /api/empleados/:id | Consultar un empleado por id |
| POST | /api/empleados | Crear un empleado |
| PUT | /api/empleados/:id | Actualizar un empleado |
| DELETE | /api/empleados/:id | Eliminar un empleado |

---

## Endpoint transaccional

**`POST /api/companias/con-empleados`**

Crea una compañía junto con varios empleados en una sola operación atómica usando Unit of Work. Si la creación de algún empleado falla, toda la operación se revierte y no se guarda ningún dato.

**Body de ejemplo:**

```json
{
  "nombre": "Tech Solutions S.A.S",
  "direccion": "Calle 45 # 10-20",
  "telefono": "3001234567",
  "empleados": [
    {
      "nombre": "Ana",
      "apellido": "Gómez",
      "correo": "ana.gomez@tech.com",
      "cargo": "Desarrolladora",
      "salario": 3500000
    },
    {
      "nombre": "Carlos",
      "apellido": "Rojas",
      "correo": "carlos.rojas@tech.com",
      "cargo": "Tester",
      "salario": 2800000
    }
  ]
}
```

---

## Instalación

```bash
# Clonar el repositorio
git clone https://github.com/alejo8933/compania-empleados-api-go.git
cd compania-empleados-api-go

# Instalar dependencias
go mod tidy
```

---

## Configuración de base de datos

Crea un archivo `.env` en la raíz del proyecto con el siguiente contenido:

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=compania_db
DB_USER=postgres
DB_PASSWORD=tu_contraseña
```

---

## Migraciones

Las migraciones se ejecutan automáticamente al iniciar el proyecto usando `AutoMigrate` de GORM. Las tablas `companias` y `empleados` se crean si no existen.

---

## Ejecución del proyecto

```bash
go run main.go
```

La API quedará disponible en `http://localhost:8080`.

---

## Pruebas con Swagger/Postman

Se pueden probar todos los endpoints usando **Postman**, **Thunder Client** o **Insomnia**.

Códigos de respuesta esperados:

| Código | Significado |
|---|---|
| 200 | Consulta exitosa |
| 201 | Recurso creado |
| 204 | Eliminación exitosa |
| 400 | Error de validación |
| 404 | Recurso no encontrado |
| 500 | Error interno del servidor |

---

## Logging

Se usa **zap** para logging estructurado. Se registran los siguientes eventos:

- Inicio de la aplicación.
- Creación de una compañía.
- Creación de un empleado.
- Inicio de una transacción.
- Confirmación de una transacción.
- Rollback de una transacción.
- Errores de base de datos.
- Errores inesperados.

---

## Uso de IA

Durante el desarrollo se usó inteligencia artificial como apoyo para:

- Entender cómo aplicar Onion Architecture en Go comparándola con ASP.NET Core.
- Identificar el ORM más adecuado para el ecosistema Go.
- Revisar conceptos de Repository Pattern y Unit of Work.
- Organizar la estructura del proyecto y la documentación.

---

## Conclusiones

Este proyecto demostró que los principios de Onion Architecture, Repository Pattern y Unit of Work son transferibles entre tecnologías. Aunque Go y ASP.NET Core son lenguajes y ecosistemas distintos, la lógica de separación de capas, el manejo de transacciones y la organización del código siguen siendo los mismos fundamentos.

Implementar estas arquitecturas en Go permitió comprender mejor la importancia de proteger la lógica de negocio de los detalles técnicos, mantener el código ordenado y garantizar la integridad de los datos mediante transacciones atómicas.
