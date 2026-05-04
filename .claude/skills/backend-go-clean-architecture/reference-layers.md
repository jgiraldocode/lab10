# Go: dependencias entre capas

```text
cmd/api/main.go
    → internal/config
    → internal/httpapi (router + handlers)
         → internal/app (services)
              → interfaces (ports) ← implementadas por internal/repository/postgres
         → internal/domain (tipos puros)
```

## Esbozo de puerto (en `internal/app` o `internal/ports`)

```go
type UserRepository interface {
    GetByID(ctx context.Context, id string) (domain.User, error)
    Create(ctx context.Context, u domain.User) error
}
```

## Esbozo de servicio

```go
type UserService struct {
    users UserRepository
}

func (s *UserService) Register(ctx context.Context, email string) error {
    // validación + políticas; delegar persistencia a s.users
    return nil
}
```

El handler solo parsea, llama `UserService`, y escribe status + JSON.
