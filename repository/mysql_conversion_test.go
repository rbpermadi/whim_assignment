package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"

	"github.com/rbpermadi/whim_assignment/app/request"
	"github.com/rbpermadi/whim_assignment/entity"
	"github.com/rbpermadi/whim_assignment/repository"

	"github.com/google/go-cmp/cmp"
)

func Test_mysqlConversion_GetConversions(t *testing.T) {
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
		want           []entity.Conversion
		wantTotal      int64
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name:      "single data",
			args:      args{context.TODO(), 10, 0},
			want:      make([]entity.Conversion, 1),
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "no args",
			args:      args{ctx: context.TODO()},
			want:      make([]entity.Conversion, 1),
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "10 data",
			args:      args{context.TODO(), 10, 0},
			want:      make([]entity.Conversion, 10),
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

			rows := sqlmock.NewRows([]string{"id", "currency_id_from", "currency_id_to", "rate", "updated_at", "created_at"})
			for _, v := range tt.want {
				rows = rows.AddRow(v.ID, v.CurrencyIDFrom, v.CurrencyIDTo, v.Rate, v.UpdatedAt, v.CreatedAt)
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

			repo := repository.NewMysqlConversion(db)
			result, total, err := repo.GetConversions(tt.args.ctx, &request.ConversionParameter{Limit: tt.args.limit, Offset: tt.args.offset})
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlConversion.GetConversions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(result) == 0 {
				return
			}

			if diff := cmp.Diff(tt.want, result); diff != "" {
				t.Errorf("mysqlConversion.GetConversions() mismatch (-want +got):\n%s", diff)
			}
			if total != tt.wantTotal {
				t.Errorf("mysqlConversion.GetConversions() got1 = %v, want %v", total, tt.wantTotal)
			}
		})
	}
}

func Test_mysqlConversion_GetConversion(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int64
	}
	sampleConversion := entity.Conversion{}
	err := faker.FakeData(&sampleConversion)
	if err != nil {
		fmt.Println(err)
	}
	tests := []struct {
		name        string
		args        args
		want        *entity.Conversion
		returnQuery error
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:    "ID 1",
			args:    args{context.TODO(), sampleConversion.ID},
			want:    &sampleConversion,
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
				rows = sqlmock.NewRows([]string{"id", "currency_id_from", "currency_id_to", "rate", "updated_at", "created_at"}).
					AddRow(tt.want.ID, tt.want.CurrencyIDFrom, tt.want.CurrencyIDTo, tt.want.Rate, tt.want.UpdatedAt, tt.want.CreatedAt)
			}

			if tt.returnQuery != nil {
				mock.ExpectQuery("^SELECT id(.+)").WillReturnError(tt.returnQuery)
			} else {
				mock.ExpectQuery("^SELECT id(.+)").WillReturnRows(rows)
			}

			repo := repository.NewMysqlConversion(db)

			gotCat, err := repo.GetConversion(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("mysqlConversion.GetConversion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, gotCat); diff != "" {
				t.Errorf("mysqlConversion.GetConversion() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_mysqlConversion_CreateConversion(t *testing.T) {
	sampleConversion := entity.Conversion{}
	err := faker.FakeData(&sampleConversion)
	if err != nil {
		fmt.Println(err)
	}
	type args struct {
		ctx      context.Context
		category *entity.Conversion
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
			args:    args{context.TODO(), &sampleConversion},
			wantErr: false,
		},
		{
			name:      "insert failed",
			args:      args{context.TODO(), &sampleConversion},
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
			prep := mock.ExpectExec("^INSERT INTO conversions(.+)")
			if tt.returnErr != nil {
				prep.WillReturnError(tt.returnErr)
			} else {
				prep.WillReturnResult(sqlmock.NewResult(2, 1))
			}

			repo := repository.NewMysqlConversion(db)
			if err := repo.CreateConversion(tt.args.ctx, tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("mysqlConversion.CreateConversion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mysqlConversion_UpdateConversion(t *testing.T) {
	sampleConversion := entity.Conversion{}
	err := faker.FakeData(&sampleConversion)
	if err != nil {
		fmt.Println(err)
	}
	type args struct {
		ctx      context.Context
		id       int64
		category *entity.Conversion
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
			args:         args{context.TODO(), sampleConversion.ID, &sampleConversion},
			rowsAffected: 1,
			wantErr:      false,
		},
		{
			name:         "update weird behaviour",
			rowsAffected: 2,
			args:         args{context.TODO(), sampleConversion.ID, &sampleConversion},
			wantErr:      true,
		},
		{
			name:      "update fail",
			args:      args{context.TODO(), sampleConversion.ID, &sampleConversion},
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

			prep := mock.ExpectExec("^UPDATE conversions(.+)")

			if tt.returnErr != nil {
				prep.WillReturnError(tt.returnErr)
			} else {
				prep.WillReturnResult(sqlmock.NewResult(2, tt.rowsAffected))
			}

			repo := repository.NewMysqlConversion(db)
			if err := repo.UpdateConversion(tt.args.ctx, tt.args.id, tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("mysqlConversion.UpdateConversion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mysqlConversion_DeleteConversion(t *testing.T) {

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

			prep := mock.ExpectExec("^DELETE FROM conversions(.+)")

			if tt.returnErr != nil {
				prep.WillReturnError(tt.returnErr)
			} else {
				prep.WillReturnResult(sqlmock.NewResult(2, tt.rowAffected))
			}

			repo := repository.NewMysqlConversion(db)
			if err := repo.DeleteConversion(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("mysqlConversion.DeleteConversion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
