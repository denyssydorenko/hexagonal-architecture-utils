package db

import (
	"context"
	"hexagonal-architexture-utils/internal/adapters/http/metrics"
	domain "hexagonal-architexture-utils/internal/domains"
	"hexagonal-architexture-utils/internal/domains/db"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

type pgxMockWithConfig struct {
	pgxmock.PgxPoolIface
}

func (pgm *pgxMockWithConfig) Config() *pgxpool.Config {
	return &pgxpool.Config{
		ConnConfig: &pgx.ConnConfig{
			Config: pgconn.Config{
				Host:     "test.host",
				Database: "testdatabase",
			},
		},
	}
}

func TestMain(m *testing.M) {
	metrics.InitMetrics()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func Test_Create(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	adapter := Adapter{
		pool: &pgxMockWithConfig{
			PgxPoolIface: mock,
		},
	}

	generatedId := uuid.NewString()

	mock.ExpectQuery("INSERT").
		WithArgs("Rebeca", "Martinez", domain.ToPointer("Spain"), int32(18)).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(generatedId))

	user := db.CreateRequest{
		Name:    "Rebeca",
		Surname: "Martinez",
		Country: domain.ToPointer("Spain"),
		Age:     int32(18),
	}
	id, err := adapter.Create(ctx, &user)
	assert.Nil(t, err)
	assert.Equal(t, id.String(), generatedId)
}

func Test_Get(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	adapter := Adapter{
		pool: &pgxMockWithConfig{
			PgxPoolIface: mock,
		},
	}

	generatedUUID := uuid.New()

	mock.ExpectQuery("SELECT").
		WithArgs(generatedUUID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "surname", "country", "age"}).
			AddRow(generatedUUID, "Mateo", "Gonzalez", nil, int32(19)))

	user, err := adapter.Get(ctx, generatedUUID)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, generatedUUID)
}

func Test_GetAll(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	adapter := Adapter{
		pool: &pgxMockWithConfig{
			PgxPoolIface: mock,
		},
	}

	mock.ExpectQuery("SELECT").
		WillReturnRows(pgxmock.NewRows([]string{"id", "name", "surname", "country", "age"}).
			AddRow(uuid.NewString(), "Mateo", "Gonzalez", nil, int32(19)).
			AddRow(uuid.NewString(), "Laura", "Dominguez", nil, int32(22)))

	resp, err := adapter.GetAll(ctx)
	assert.Nil(t, err)
	assert.Equal(t, len(*resp), 2)
}

func Test_Update(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	adapter := Adapter{
		pool: &pgxMockWithConfig{
			PgxPoolIface: mock,
		},
	}

	generatedUUID := uuid.New()

	mock.ExpectExec("UPDATE").
		WithArgs(generatedUUID, "Rebeca", "Martinez").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	id, err := adapter.Update(ctx, generatedUUID, "Rebeca", "Martinez")
	assert.Nil(t, err)
	assert.Equal(t, *id, generatedUUID)
}

func Test_Delete(t *testing.T) {
	ctx := context.Background()

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	adapter := Adapter{
		pool: &pgxMockWithConfig{
			PgxPoolIface: mock,
		},
	}

	generatedUUID := uuid.New()

	mock.ExpectExec("DELETE").
		WithArgs(generatedUUID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	deleted, err := adapter.Delete(ctx, generatedUUID)
	assert.Nil(t, err)
	assert.Equal(t, *deleted, true)
}
