package api

import (
	"context"
	database "hexagonal-architexture-utils/internal/adapters/db"
	"hexagonal-architexture-utils/internal/domains/db"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_DB_Create(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mock := database.NewMockDBPort(ctrl)

	id := uuid.New()
	user := &db.CreateRequest{
		Name:    "Test",
		Surname: "Test",
		Age:     20,
	}
	mock.EXPECT().Create(ctx, user).Return(&id, nil)

	api := NewApplication(mock)
	resp, err := api.DBCreate(ctx, user)
	assert.Nil(t, err)
	assert.Equal(t, resp.Id, id)
}

func Test_DB_Get(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mock := database.NewMockDBPort(ctrl)

	id := uuid.New()
	user := &db.User{
		ID:      id,
		Name:    "Test",
		Surname: "Test",
		Age:     20,
	}
	mock.EXPECT().Get(ctx, id).Return(user, nil)

	api := NewApplication(mock)
	resp, err := api.DBGet(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, resp.Name, "Test")
}

func Test_DB_GetAll(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mock := database.NewMockDBPort(ctrl)

	id1 := uuid.New()
	id2 := uuid.New()
	users := []db.User{
		{
			ID:      id1,
			Name:    "Test",
			Surname: "Test",
			Age:     20,
		},
		{
			ID:      id2,
			Name:    "Test2",
			Surname: "Test2",
			Age:     22,
		},
	}
	mock.EXPECT().GetAll(ctx).Return(&users, nil)

	api := NewApplication(mock)
	resp, err := api.DBGetAll(ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(*resp), 2)
}

func Test_DB_Update(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mock := database.NewMockDBPort(ctrl)

	id := uuid.New()
	mock.EXPECT().Update(ctx, id, "Test", "test").Return(&id, nil)

	api := NewApplication(mock)

	req := db.UpdateRequest{
		Id:      id,
		Name:    "Test",
		Surname: "test",
	}
	resp, err := api.DBUpdate(ctx, &req)
	assert.Nil(t, err)
	assert.Equal(t, resp.Id, id)
}

func Test_DB_Delete(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mock := database.NewMockDBPort(ctrl)

	id := uuid.New()
	boolResp := true
	mock.EXPECT().Delete(ctx, id).Return(&boolResp, nil)

	api := NewApplication(mock)

	resp, err := api.DBDelete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, *resp, true)
}
