package location_repo

import (
	"context"
	"final-project/internal/location/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type locationRepo struct {
	pgxPool *pgxpool.Pool
}

func (l locationRepo) UpdateLocation(ctx context.Context, lat, lng float64, driverID string) error {
	_, err := l.pgxPool.Exec(ctx, `WITH upsert AS (
			 UPDATE driver_location SET lat=$1, lng=$2 
			 WHERE driver_id=$3 
			 RETURNING *
		)
		INSERT INTO driver_location (lat, lng, driver_id) 
		SELECT $1, $2, $3 
		WHERE NOT EXISTS (SELECT * FROM upsert)`, lat, lng, driverID)
	if err != nil {
		return err
	}

	return nil
}

func (l locationRepo) GetNearbyDrivers(ctx context.Context, lat, lng, radius float64) error {
	panic("implement me")
}

func New(pgxPool *pgxpool.Pool) repository.Location {
	return &locationRepo{pgxPool: pgxPool}
}
