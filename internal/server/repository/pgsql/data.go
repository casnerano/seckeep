package pgsql

import (
	"context"
	"errors"
	"time"

	"github.com/casnerano/seckeep/internal/pkg"
	"github.com/casnerano/seckeep/internal/server/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DataRepository структура репозитория работы с записями секретных данных.
type DataRepository struct {
	pgxpool *pgxpool.Pool
}

// NewDataRepository конструктор.
func NewDataRepository(pgxpool *pgxpool.Pool) repository.Data {
	return &DataRepository{pgxpool}
}

// Add добавляет запись.
func (d DataRepository) Add(ctx context.Context, data pkg.Data) (*pkg.Data, error) {
	err := d.pgxpool.QueryRow(
		ctx,
		"insert into data(user_uuid, type, value, created_at, version) values($1, $2, $3, $4, $5) returning uuid",
		data.UserUUID,
		data.Type,
		data.Value,
		data.CreatedAt.UTC(),
		data.Version.UTC(),
	).Scan(
		&data.UUID,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// FindByUUID ищет запись по UUID.
func (d DataRepository) FindByUUID(ctx context.Context, userUUID string, uuid string) (*pkg.Data, error) {
	data := pkg.Data{UUID: uuid}
	err := d.pgxpool.QueryRow(
		ctx,
		"select user_uuid, type, value, created_at, version from data where user_uuid = $1 and uuid = $2",
		userUUID,
		uuid,
	).Scan(
		&data.UserUUID,
		&data.Type,
		&data.Value,
		&data.CreatedAt,
		&data.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrNotFound
		}
		return nil, err
	}

	return &data, nil
}

// FindByUserUUID ищет запись по UUID пользователя.
func (d DataRepository) FindByUserUUID(ctx context.Context, userUUID string) ([]*pkg.Data, error) {
	data := make([]*pkg.Data, 0)

	rows, err := d.pgxpool.Query(
		ctx,
		"select uuid, type, value, created_at, version from data where user_uuid = $1",
		userUUID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		datum := &pkg.Data{}
		err = rows.Scan(
			&datum.UUID,
			&datum.Type,
			&datum.Value,
			&datum.CreatedAt,
			&datum.Version,
		)
		if err == nil {
			data = append(data, datum)
		}
	}

	return data, nil
}

// Update обновляет запись.
func (d DataRepository) Update(ctx context.Context, userUUID string, uuid string, value []byte, version time.Time) (*pkg.Data, error) {
	data := &pkg.Data{
		UUID:    uuid,
		Value:   value,
		Version: version,
	}

	err := d.pgxpool.QueryRow(
		ctx,
		"update data set value = $1, version = $2 where user_uuid = $3 and uuid = $4 returning user_uuid, type, created_at",
		value,
		data.Version.UTC(),
		userUUID,
		uuid,
	).Scan(
		&data.UserUUID,
		&data.Type,
		&data.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repository.ErrNotFound
		}
		return nil, err
	}

	return data, nil
}

// Delete удаляет запись.
func (d DataRepository) Delete(ctx context.Context, userUUID string, uuid string) error {
	res, err := d.pgxpool.Exec(
		ctx,
		"delete from data where user_uuid = $1 and uuid = $2",
		userUUID,
		uuid,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() > 0 {
		return nil
	}

	return repository.ErrNotFound
}
