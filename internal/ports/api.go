package ports

import (
	"context"
	"hexagonal-architexture-utils/internal/domains/db"

	"github.com/google/uuid"
)

type ApiPort interface {
	IsHealthy(ctx context.Context) bool

	// DB
	DBCreate(ctx context.Context, user *db.CreateRequest) (*db.ID, error)
	DBGet(ctx context.Context, id uuid.UUID) (*db.User, error)
	DBGetAll(ctx context.Context) (*[]db.User, error)
	DBUpdate(ctx context.Context, req *db.UpdateRequest) (*db.ID, error)
	DBDelete(ctx context.Context, id uuid.UUID) (*bool, error)
}
