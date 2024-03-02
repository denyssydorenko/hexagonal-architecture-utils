package api

import (
	"context"
	"hexagonal-architexture-utils/internal/domains/db"
	"hexagonal-architexture-utils/internal/ports"

	"github.com/google/uuid"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{db: db}
}

func (app *Application) IsHealthy(ctx context.Context) bool {
	return app.db.Health(ctx)
}

func (app *Application) DBCreate(ctx context.Context, user *db.CreateRequest) (*db.ID, error) {
	id, err := app.db.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &db.ID{Id: *id}, nil
}

func (app *Application) DBGet(ctx context.Context, id uuid.UUID) (*db.User, error) {
	return app.db.Get(ctx, id)
}

func (app *Application) DBGetAll(ctx context.Context) (*[]db.User, error) {
	return app.db.GetAll(ctx)
}

func (app *Application) DBUpdate(ctx context.Context, req *db.UpdateRequest) (*db.ID, error) {
	resp, err := app.db.Update(ctx, req.Id, req.Name, req.Surname)
	if err != nil {
		return nil, err
	}
	return &db.ID{Id: *resp}, nil
}

func (app *Application) DBDelete(ctx context.Context, id uuid.UUID) (*bool, error) {
	return app.db.Delete(ctx, id)
}
