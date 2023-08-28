package repository

import (
	"avitoTech/config"
	models "avitoTech/internal"
	"avitoTech/utils"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"
)

type postgres struct {
	pool   *pgxpool.Pool
	logger *log.Logger
}

func NewPostgres(ctx context.Context, config *config.Config) (*postgres, error) {
	logger := log.New(os.Stderr, "Postgres: ", log.LstdFlags|log.Lshortfile)
	pool, err := utils.GetPool(ctx, config)
	if err != nil {
		logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return nil, &respond
	}
	return &postgres{
		pool:   pool,
		logger: logger,
	}, nil
}

func (p *postgres) CreateSegment(ctx context.Context, segmentName string) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}

	_, err = conn.Exec(ctx, `insert into segmentsEnum (segmentName) values ($1)`, segmentName)
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}

	return nil
}

func (p *postgres) DeleteSegment(ctx context.Context, segmentName string) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}

	_, err = conn.Exec(ctx, `delete from segmentsEnum where segmentName = $1`, segmentName)
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}
	return nil
}

func (p *postgres) UnSubscribeFromSegment(ctx context.Context, segmentName string) error {
	//var segmentId int

	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}

	_, err = conn.Exec(ctx, `DELETE FROM mainTable 
       WHERE segmentId IN (SELECT segmentId FROM segmentsEnum WHERE segmentName = $1)`, segmentName)
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}

	return nil
}

func (p *postgres) SubscribeUserOnSegment(ctx context.Context, userId int, segmentName string) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}

	var t int
	err = conn.QueryRow(ctx, `WITH segment AS (
    SELECT segmentId FROM segmentsEnum WHERE segmentName = $1
)
INSERT INTO mainTable (userId, segmentId)
SELECT $2, segmentId
FROM segment
RETURNING segmentId;`, segmentName, userId).Scan(&t)

	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		if errors.Is(err, pgx.ErrNoRows) {
			respond.Message = "this segment does not exist"
		}
		return &respond
	}

	return nil
}

func (p *postgres) SubscribeUserOnSegmentWithTimeout(ctx context.Context, userId int, segmentName string, timeout time.Time) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}
	var t int
	err = conn.QueryRow(ctx, `WITH segment AS (
    SELECT segmentId FROM segmentsEnum WHERE segmentName = $1
)
INSERT INTO mainTable (userId, segmentId, timeEnd)
SELECT $2, segmentId, $3
FROM segment
RETURNING segmentId;`, segmentName, userId, timeout).Scan(&t)

	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		if errors.Is(err, pgx.ErrNoRows) {
			respond.Message = "this segment does not exist"
		}
		return &respond
	}

	return nil
}

func (p *postgres) UnsubscribeUser(ctx context.Context, userId int, segmentName string) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}

		return &respond
	}

	_, err = conn.Exec(ctx, `WITH segment AS (
    SELECT segmentId FROM segmentsEnum WHERE segmentName = $2
) 
DELETE FROM mainTable 
WHERE userId = $1 AND segmentId IN (SELECT segmentId FROM segment);`,
		userId, segmentName)
	if err != nil {
		p.logger.Print(err)
		respond := models.Respond{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &respond
	}
	return nil
}
