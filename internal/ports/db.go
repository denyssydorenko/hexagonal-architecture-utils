package ports

import (
	"context"
	"hexagonal-architexture-utils/internal/domains/db"

	"github.com/google/uuid"
)

type DBPort interface {
	Health(ctx context.Context) bool
	Close()

	Create(ctx context.Context, user *db.CreateRequest) (*uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*db.User, error)
	GetAll(ctx context.Context) (*[]db.User, error)
	Update(ctx context.Context, id uuid.UUID, name string, surname string) (*uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) (*bool, error)
}
