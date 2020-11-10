package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"

	"github.com/google/go-cmp/cmp"
	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/repository"
)

func Test_mysqlCurrency_GetCurrencies(t *testing.T) {
	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}

	tests := []struct {
		name           string
		args           args
		countErrQuery  error
		selectErrQuery error
		want           []entity.Currency
		wantTotal      int64
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name:      "single data",
			args:      args{context.TODO(), 10, 0},
			want:      make([]entity.Currency, 1),
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "no args",
			args:      args{ctx: context.TODO()},
			want:      make([]entity.Currency, 1),
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "10 data",
			args:      args{context.TODO(), 10, 0},
			want:      make([]entity.Currency, 10),
			wantTotal: 10,
			wantErr:   false,
		},
		{
			name:          "error count",
			args:          args{context.TODO(), 10, 0},
			countErrQuery: errors.New("error count"),
			wantErr:       true,
		},
		{
			name:           "error select",
			args:           args{context.TODO(), 10, 0},
			selectErrQuery: errors.New("error select"),
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.want {
				err := faker.FakeData(&tt.want[i])
				if err != nil {
					fmt.Println(err)
				}
			}

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"id", "name", "updated_at", "created_at"})
			for _, v := range tt.want {
				rows = rows.AddRow(v.ID, v.Name, v.UpdatedAt, v.CreatedAt)
			}

			rowCount := sqlmock.NewRows([]string{"total"}).AddRow(len(tt.want))
			// Just mock regex query with respective flow query,
			// in this case Fetch function call select count then select data
			if tt.countErrQuery != nil {
				mock.ExpectQuery("^SELECT COUNT(.+)").WillReturnError(tt.countErrQuery)
			} else {
				mock.ExpectQuery("^SELECT COUNT(.+)").WillReturnRows(rowCount)
			}

			if tt.selectErrQuery != nil {
				mock.ExpectQuery("^SELECT id(.+)").WillReturnError(tt.selectErrQuery)
			} else {
				mock.ExpectQuery("^SELECT id(.+)").WillReturnRows(rows)
			}

			repo := repository.NewMysqlCurrency(db)
			result, total, err := repo.GetCurrencies(tt.args.ctx, &request.CurrencyParameter{Limit: tt.args.limit, Offset: tt.args.offset})
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlCurrency.GetCurrencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(result) == 0 {
				return
			}

			if diff := cmp.Diff(tt.want, result); diff != "" {
				t.Errorf("mysqlCurrency.GetCurrencies() mismatch (-want +got):\n%s", diff)
			}
			if total != tt.wantTotal {
				t.Errorf("mysqlCurrency.GetCurrencies() got1 = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func Test_mysqlCurrency_GetCurrency(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	sampleCurrency := entity.Currency{}
	err := faker.FakeData(&sampleCurrency)
	if err != nil {
		fmt.Println(err)
	}
	tests := []struct {
		name        string
		args        args
		want        *entity.Currency
		returnQuery error
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:    "ID 1",
			args:    args{context.TODO(), sampleCurrency.ID},
			want:    &sampleCurrency,
			wantErr: false,
		},
		{
			name:        "not found error",
			args:        args{context.TODO(), 1},
			returnQuery: sql.ErrNoRows,
			wantErr:     true,
		},
		{
			name:        "unknown error",
			args:        args{context.TODO(), 1},
			returnQuery: fmt.Errorf("sql error: %s", "unknown error"),
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rows *sqlmock.Rows
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			if tt.want != nil {
				rows = sqlmock.NewRows([]string{"id", "name", "updated_at", "created_at"}).
					AddRow(tt.want.ID, tt.want.Name, tt.want.UpdatedAt, tt.want.CreatedAt)
			}

			if tt.returnQuery != nil {
				mock.ExpectQuery("^SELECT id(.+)").WillReturnError(tt.returnQuery)
			} else {
				mock.ExpectQuery("^SELECT id(.+)").WillReturnRows(rows)
			}

			repo := repository.NewMysqlCurrency(db)

			gotCat, err := repo.GetCurrency(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlCurrency.GetCurrency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, gotCat); diff != "" {
				t.Errorf("mysqlCurrency.GetCurrency() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_mysqlCurrency_CreateCurrency(t *testing.T) {
	sampleCurrency := entity.Currency{}
	err := faker.FakeData(&sampleCurrency)
	if err != nil {
		fmt.Println(err)
	}
	type args struct {
		ctx      context.Context
		currency *entity.Currency
	}
	tests := []struct {
		name      string
		args      args
		returnErr error
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name:    "insert ok",
			args:    args{context.TODO(), &sampleCurrency},
			wantErr: false,
		},
		{
			name:      "insert failed",
			args:      args{context.TODO(), &sampleCurrency},
			returnErr: errors.New("fail insert"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			prep := mock.ExpectExec("^INSERT INTO currencies(.+)")
			if tt.returnErr != nil {
				prep.WillReturnError(tt.returnErr)
			} else {
				prep.WillReturnResult(sqlmock.NewResult(2, 1))
			}

			repo := repository.NewMysqlCurrency(db)
			if err := repo.CreateCurrency(tt.args.ctx, tt.args.currency); (err != nil) != tt.wantErr {
				t.Errorf("mysqlCurrency.CreateCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mysqlCurrency_UpdateCurrency(t *testing.T) {
	sampleCurrency := entity.Currency{}
	err := faker.FakeData(&sampleCurrency)
	if err != nil {
		fmt.Println(err)
	}
	type args struct {
		ctx      context.Context
		id       int64
		currency *entity.Currency
	}
	tests := []struct {
		name         string
		args         args
		returnErr    error
		rowsAffected int64
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			name:         "update ok",
			args:         args{context.TODO(), sampleCurrency.ID, &sampleCurrency},
			rowsAffected: 1,
			wantErr:      false,
		},
		{
			name:         "update weird behaviour",
			rowsAffected: 2,
			args:         args{context.TODO(), sampleCurrency.ID, &sampleCurrency},
			wantErr:      true,
		},
		{
			name:      "update fail",
			args:      args{context.TODO(), sampleCurrency.ID, &sampleCurrency},
			returnErr: errors.New("update fail"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			prep := mock.ExpectExec("^UPDATE currencies(.+)")

			if tt.returnErr != nil {
				prep.WillReturnError(tt.returnErr)
			} else {
				prep.WillReturnResult(sqlmock.NewResult(2, tt.rowsAffected))
			}

			repo := repository.NewMysqlCurrency(db)
			if err := repo.UpdateCurrency(tt.args.ctx, tt.args.id, tt.args.currency); (err != nil) != tt.wantErr {
				t.Errorf("mysqlCurrency.UpdateCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mysqlCurrency_DeleteCurrency(t *testing.T) {

	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name        string
		args        args
		returnErr   error
		rowAffected int64
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:        "delete ok",
			args:        args{context.TODO(), 1},
			rowAffected: 1,
			wantErr:     false,
		},
		{
			name:        "delete weird",
			args:        args{context.TODO(), 1},
			rowAffected: 2,
			wantErr:     true,
		},
		{
			name:      "delete failed",
			args:      args{context.TODO(), 1},
			returnErr: errors.New("fail delete"),
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			prep := mock.ExpectExec("^DELETE FROM currencies(.+)")

			if tt.returnErr != nil {
				prep.WillReturnError(tt.returnErr)
			} else {
				prep.WillReturnResult(sqlmock.NewResult(2, tt.rowAffected))
			}

			repo := repository.NewMysqlCurrency(db)
			if err := repo.DeleteCurrency(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("mysqlCurrency.DeleteCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
