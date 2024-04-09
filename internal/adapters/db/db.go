package db

import (
	"context"
	"fmt"
	"hexagonal-architexture-utils/config"
	"hexagonal-architexture-utils/internal/adapters/http/metrics"
	domain "hexagonal-architexture-utils/internal/domains"
	"hexagonal-architexture-utils/internal/domains/db"
	"regexp"
	"strconv"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Adapter struct {
	pool PgxPoolIface
}

type PgxPoolIface interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Config() *pgxpool.Config
	Ping(context.Context) error
	Close()
}

func NewAdapter(ctx context.Context, conf *config.Config) (*Adapter, error) {
	connString := getConnectionString(conf)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool config: %w", err)
	}
	config.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithIncludeQueryParameters())

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to postgres [%v]", err)
	}
	return &Adapter{
		pool: pool,
	}, nil
}

func getConnectionString(conf *config.Config) string {
	return fmt.Sprintf("host='%v'port='%v'dbname='%v'user='%v'password='%v'pool_min_conns='%v'",
		conf.DBHost(),
		uint16(conf.DBPort()),
		conf.DBName(),
		conf.DBUser(),
		conf.DBPassword(),
		conf.MinOpenConns(),
	)
}

func (ad *Adapter) Health(ctx context.Context) bool {
	var one int
	err := ad.pool.QueryRow(ctx, `SELECT 1`).Scan(&one)
	if err != nil {
		return false
	}
	return true
}

func (ad *Adapter) Close() {
	ad.pool.Close()
}

func (ad *Adapter) Create(ctx context.Context, user *db.CreateRequest) (*uuid.UUID, error) {
	query := `INSERT INTO users (name, surname, country, age) VALUES ($1, $2, $3, $4) RETURNING id`

	start := time.Now()

	var id uuid.UUID
	err := ad.pool.QueryRow(ctx, query, user.Name, user.Surname, user.Country, user.Age).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("was not possible to insert the user: %+v", err)
	}

	metrics.DBDuration.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "Create", strconv.FormatBool(err == nil)).Observe(time.Since(start).Seconds())
	metrics.DBQueries.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "Create", strconv.FormatBool(err == nil)).Inc()

	return &id, nil
}

func (ad *Adapter) Get(ctx context.Context, id uuid.UUID) (*db.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	start := time.Now()

	rows, err := ad.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("was not possible to get the users: %w", err)
	}
	defer rows.Close()

	metrics.DBDuration.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "Get", strconv.FormatBool(err == nil)).Observe(time.Since(start).Seconds())
	metrics.DBQueries.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "Get", strconv.FormatBool(err == nil)).Inc()

	var user db.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Country, &user.Age); err != nil {
			return nil, fmt.Errorf("was not possible to parse db row to user: %w", err)
		}
	}

	if IsValidUUID(user.ID.String()) {
		return &user, nil
	}

	return nil, nil
}

func (ad *Adapter) GetAll(ctx context.Context) (*[]db.User, error) {
	query := `SELECT * FROM users`

	start := time.Now()

	rows, err := ad.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("was not possible to get the users: %w", err)
	}
	defer rows.Close()

	metrics.DBDuration.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "GetAll", strconv.FormatBool(err == nil)).Observe(time.Since(start).Seconds())
	metrics.DBQueries.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "GetAll", strconv.FormatBool(err == nil)).Inc()

	var users []db.User
	for rows.Next() {
		var user db.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Country, &user.Age); err != nil {
			return nil, fmt.Errorf("was not possible to parse db row to user: %w", err)
		}
		users = append(users, user)
	}

	if users != nil {
		return &users, nil
	}

	return nil, nil
}

func (ad *Adapter) Update(ctx context.Context, id uuid.UUID, name string, surname string) (*uuid.UUID, error) {
	query := `UPDATE users SET name = $2, surname = $3 WHERE id = $1 RETURNING id`

	start := time.Now()

	_, err := ad.pool.Exec(ctx, query, id, name, surname)
	if err != nil {
		return nil, fmt.Errorf("was not possible to update the user: %w", err)
	}

	metrics.DBDuration.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "Update", strconv.FormatBool(err == nil)).Observe(time.Since(start).Seconds())
	metrics.DBQueries.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "Update", strconv.FormatBool(err == nil)).Inc()

	return &id, nil
}

func (ad *Adapter) Delete(ctx context.Context, id uuid.UUID) (*bool, error) {
	query := `DELETE FROM users WHERE id = $1`

	start := time.Now()

	_, err := ad.pool.Exec(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("was not possible to update the user: %w", err)
	}

	metrics.DBDuration.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "Delete", strconv.FormatBool(err == nil)).Observe(time.Since(start).Seconds())
	metrics.DBQueries.WithLabelValues(ad.pool.Config().ConnConfig.Host, ad.pool.Config().ConnConfig.Database, "Delete", strconv.FormatBool(err == nil)).Inc()

	return domain.ToPointer(true), nil
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
