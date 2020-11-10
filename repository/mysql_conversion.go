package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
)

type mysqlConversion struct {
	db *sql.DB
}

type ConversionRepo interface {
	CreateConversion(ctx context.Context, ec *entity.Conversion) error
	UpdateConversion(ctx context.Context, id int64, ec *entity.Conversion) error
	DeleteConversion(ctx context.Context, id int64) error
	GetConversion(ctx context.Context, id int64) (*entity.Conversion, error)
	GetConversions(ctx context.Context, p *request.ConversionParameter) ([]entity.Conversion, int64, error)
}

//NewMysqlConversion is a function to create implementation of mysql Conversion repository
func NewMysqlConversion(db *sql.DB) ConversionRepo {
	return &mysqlConversion{db}
}

func (t *mysqlConversion) fetch(ctx context.Context, query string) ([]entity.Conversion, error) {
	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]entity.Conversion, 0)
	for rows.Next() {
		cat := entity.Conversion{}
		err = rows.Scan(
			&cat.ID,
			&cat.CurrencyIDFrom,
			&cat.CurrencyIDTo,
			&cat.Rate,
			&cat.UpdatedAt,
			&cat.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		result = append(result, cat)
	}

	return result, nil
}

func (t *mysqlConversion) GetConversion(ctx context.Context, id int64) (*entity.Conversion, error) {
	query := `SELECT id, currency_id_from, currency_id_to, rate, updated_at, created_at
						  FROM conversions WHERE id = %d`

	list, err := t.fetch(ctx, buildQuery(query, id))
	if err == sql.ErrNoRows || len(list) == 0 {
		return nil, fmt.Errorf("Not Found")
	}

	return &list[0], nil
}

func (t *mysqlConversion) GetConversions(ctx context.Context, p *request.ConversionParameter) ([]entity.Conversion, int64, error) {
	var result []entity.Conversion
	var total int64
	var queryString string

	if p.CurrencyIDFrom != 0 && p.CurrencyIDTo != 0 {
		queryCount := "SELECT COUNT(id) FROM conversions WHERE (currency_id_from = %d AND currency_id_to = %d) OR (currency_id_from = %d AND currency_id_to = %d)"
		err := t.db.QueryRowContext(ctx, buildQuery(queryCount, p.CurrencyIDFrom, p.CurrencyIDTo, p.CurrencyIDTo, p.CurrencyIDFrom)).Scan(&total)
		if err != nil {
			return nil, 0, err
		}

		query := `SELECT
								id, currency_id_from, currency_id_to, rate, updated_at, created_at
							FROM
								conversions
							WHERE (currency_id_from = %d AND currency_id_to = %d) OR (currency_id_from = %d AND currency_id_to = %d)
							LIMIT %d, %d `
		queryString = buildQuery(query, p.CurrencyIDFrom, p.CurrencyIDTo, p.CurrencyIDTo, p.CurrencyIDFrom, p.Offset, p.Limit)
	} else {
		err := t.db.QueryRowContext(ctx, "SELECT COUNT(id) FROM conversions").Scan(&total)
		if err != nil {
			return nil, 0, err
		}

		query := `SELECT id, currency_id_from, currency_id_to, rate, updated_at, created_at FROM conversions LIMIT %d, %d `
		queryString = buildQuery(query, p.Offset, p.Limit)
	}

	result, err := t.fetch(ctx, queryString)
	if err != nil {
		return nil, 0, err
	}

	return result, total, err
}

func (t *mysqlConversion) CreateConversion(ctx context.Context, conversion *entity.Conversion) error {
	query := `INSERT INTO conversions (currency_id_from, currency_id_to, rate, updated_at, created_at) VALUES (%d, %d, %f, %q, %q)`
	res, err := t.db.ExecContext(ctx,
		buildQuery(query,
			conversion.CurrencyIDFrom,
			conversion.CurrencyIDTo,
			conversion.Rate,
			sqlTime(conversion.UpdatedAt),
			sqlTime(conversion.CreatedAt)),
	)

	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	conversion.ID = lastID
	return nil
}

func (t *mysqlConversion) UpdateConversion(ctx context.Context, id int64, Conversion *entity.Conversion) error {
	query := `UPDATE conversions set rate=%f, updated_at=%q WHERE ID = %d`

	res, err := t.db.ExecContext(ctx, buildQuery(query, Conversion.Rate, sqlTime(Conversion.UpdatedAt), id))
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect == 0 {
		err = fmt.Errorf("Not Found")

		return err
	}

	if affect != 1 {
		err = fmt.Errorf("weird  behaviour. total affected: %d", affect)

		return err
	}

	return nil
}

func (t *mysqlConversion) DeleteConversion(ctx context.Context, id int64) error {
	query := "DELETE FROM conversions WHERE id = %d"

	res, err := t.db.ExecContext(ctx, buildQuery(query, id))
	if err != nil {

		return err
	}
	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAfected != 1 {
		err = fmt.Errorf("weird behaviour. total affected: %d", rowsAfected)
		return err
	}

	return nil
}
