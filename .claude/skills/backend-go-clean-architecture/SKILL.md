---
name: backend-go-clean-architecture
description: >-
  Implements Go HTTP APIs with clear separation between transport, application
  logic, persistence, and configuration (layered / hexagonal style). Use when the
  user builds backends in Go, REST or JSON APIs, PostgreSQL/SQL, migrations,
  repositories, or asks for scalable project layout.
---

# Backend: Go con arquitectura en capas

## Objetivo

Separación **explícita** entre: entrada HTTP (transporte), casos de uso (app), dominio, acceso a datos, e infraestructura (config, logging, DB). Dependencias apuntan **hacia dentro** (dominio y puertos no importan implementaciones concretas).

## Layout de proyecto recomendado

```
cmd/
  api/
    main.go              # wiring: config, db, router, server listen
internal/
  config/
    config.go            # env + validación
  httpapi/               # o server/, transport/
    router.go
    middleware.go
    handlers/            # un archivo por recurso o por agregado
      users.go
  app/                   # application / use cases
    users_service.go     # orquesta repositorios + reglas de aplicación
  domain/                # entidades y errores de dominio (sin imports de sql, echo, etc.)
    user.go
    errors.go
  repository/            # implementaciones de persistencia
    postgres/
      users_repo.go
      db.go              # pool, ping
  platform/              # opcional: email, colas, clock
pkg/                     # solo código realmente genérico y estable
  httputil/
  validation/
migrations/              # SQL (goose, golang-migrate, atlas) o embed
go.mod
```

- **`cmd/`**: puntos de entrada delgados; no lógica de negocio.
- **`internal/`**: no importable desde otros módulos; aquí vive el producto.
- **`pkg/`**: usar con moderación; si solo un servicio lo usa, mantenerlo en `internal`.

## Capas y responsabilidades

| Capa | Contenido | No debe |
|------|-------------|---------|
| **handlers** (`httpapi/handlers`) | Decode JSON, códigos HTTP, path/query params, llamar al servicio de aplicación | SQL, reglas de negocio pesadas |
| **app** | Casos de uso, transacciones, coordinar varios repos | Detalles de `net/http` o frameworks |
| **domain** | Structs de negocio, invariantes, errores tipados | Drivers de DB, tags JSON de API si se quiere desacoplar al máximo |
| **repository** | Queries, mapeo fila↔dominio, timeouts de contexto | Lógica de autorización de producto |

Los **handlers** dependen de **interfaces** definidas al lado del consumidor (en `app` o en un paquete `ports`): por ejemplo `UserReader`, `UserWriter`. `repository/postgres` implementa esas interfaces.

## HTTP y contexto

- Usar `context.Context` en **toda** cadena pública (handlers → app → repo).
- Timeouts en servidor HTTP y en consultas DB (`QueryContext`, etc.).
- Respuestas JSON coherentes: struct de error `{ "error": { "code", "message" } }` alineada con el proyecto.

## Base de datos

- **Una** función en `repository/postgres` (o equivalente) que abre el pool (`pgxpool` o `database/sql` + driver).
- Migraciones versionadas en `migrations/`; nunca esquema solo desde código salvo convención explícita del equipo.
- Transacciones en la capa **app** cuando un caso de uso toca varias tablas: pasar `tx` vía interfaz o método `WithinTx` según patrón elegido; ser consistente.

## Configuración

- Leer de variables de entorno en `internal/config`; validar al arranque (campos obligatorios, URLs, puertos).
- No leer `os.Getenv` disperso en handlers.

## Errores

- Errores de dominio como valores o tipos (`var ErrNotFound = ...`) mapeados a status HTTP en handlers o en un middleware central.
- No filtrar detalles internos (SQL, paths) al cliente en producción.

## Testing

- **Domain / app**: tests con tablas y mocks de interfaces de repo.
- **Repository**: tests de integración con contenedor o `testdata` cuando el proyecto lo tenga.
- Handlers: tests HTTP con `httptest` y cuerpo JSON.

## Checklist al crear o extender un endpoint

- [ ] Contrato (ruta, método, DTO request/response) claro.
- [ ] Lógica en `app`, no en el handler.
- [ ] Repo solo acceso a datos; sin reglas de producto duplicadas.
- [ ] Context y errores propagados correctamente.
- [ ] Migración si el esquema cambia.

## Referencia ampliada

Para diagrama de dependencias y ejemplo de interfaces, ver [reference-layers.md](reference-layers.md).
