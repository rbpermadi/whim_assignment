package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
)

type mysqlCurrency struct {
	db *sql.DB
}

type CurrencyRepo interface {
	CreateCurrency(ctx context.Context, ec *entity.Currency) error
	UpdateCurrency(ctx context.Context, id int64, ec *entity.Currency) error
	DeleteCurrency(ctx context.Context, id int64) error
	GetCurrency(ctx context.Context, id int64) (*entity.Currency, error)
	GetCurrencies(ctx context.Context, p *request.CurrencyParameter) ([]entity.Currency, int64, error)
}

//NewMysqlCurrency is a function to create implementation of mysql Currency repository
func NewMysqlCurrency(db *sql.DB) CurrencyRepo {
	return &mysqlCurrency{db}
}

func (t *mysqlCurrency) fetch(ctx context.Context, query string) ([]entity.Currency, error) {
	rows, err := t.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make([]entity.Currency, 0)
	for rows.Next() {
		cat := entity.Currency{}
		err = rows.Scan(
			&cat.ID,
			&cat.Name,
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

func (t *mysqlCurrency) GetCurrency(ctx context.Context, id int64) (*entity.Currency, error) {
	query := `SELECT id, name, updated_at, created_at
						  FROM currencies WHERE id = %d`

	list, err := t.fetch(ctx, buildQuery(query, id))
	if err == sql.ErrNoRows || len(list) == 0 {
		return nil, fmt.Errorf("Not Found")
	}

	return &list[0], nil
}

func (t *mysqlCurrency) GetCurrencies(ctx context.Context, p *request.CurrencyParameter) ([]entity.Currency, int64, error) {
	var result []entity.Currency
	var total int64
	var queryString string

	err := t.db.QueryRowContext(ctx, "SELECT COUNT(id) FROM currencies").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if p.Query != "" {
		query := `SELECT id, name, updated_at, created_at FROM currencies WHERE name LIKE '%%%s%%' LIMIT %d, %d `
		queryString = buildQuery(query, p.Query, p.Offset, p.Limit)
		fmt.Println(queryString)
	} else {
		query := `SELECT id, name, updated_at, created_at FROM currencies LIMIT %d, %d `
		queryString = buildQuery(query, p.Offset, p.Limit)
	}

	result, err = t.fetch(ctx, queryString)
	if err != nil {
		return nil, 0, err
	}

	return result, total, err
}

func (t *mysqlCurrency) CreateCurrency(ctx context.Context, Currency *entity.Currency) error {
	query := `INSERT INTO currencies (name, updated_at, created_at) VALUES (%q, %q, %q)`

	res, err := t.db.ExecContext(ctx,
		buildQuery(query,
			Currency.Name,
			sqlTime(Currency.UpdatedAt),
			sqlTime(Currency.CreatedAt)),
	)
	if err != nil {

		return err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	Currency.ID = lastID
	return nil
}

func (t *mysqlCurrency) UpdateCurrency(ctx context.Context, id int64, Currency *entity.Currency) error {
	query := `UPDATE currencies set name=%q, updated_at=%q WHERE ID = %d`

	res, err := t.db.ExecContext(ctx, buildQuery(query, Currency.Name, sqlTime(Currency.UpdatedAt), id))
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

func (t *mysqlCurrency) DeleteCurrency(ctx context.Context, id int64) error {
	query := "DELETE FROM currencies WHERE id = %d"

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
