package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"test/pkg/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertUrl(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &RepoPostgre{
		db: db,
	}

	testCases := map[string]struct {
		short      string
		long       string
		resultErr  error
		mockResult sql.Result
	}{
		"OK": {
			short:      "shrt",
			long:       "longg",
			resultErr:  nil,
			mockResult: sqlmock.NewResult(0, 1),
		},
		"Bad insert": {
			short:      "",
			long:       "lg",
			resultErr:  fmt.Errorf("repo err"),
			mockResult: sqlmock.NewErrorResult(fmt.Errorf("repo err")),
		},
	}

	selectRow := "INSERT INTO url(long, short) VALUES ($1, $2)"
	for _, curr := range testCases {
		if curr.resultErr == nil {
			mock.ExpectExec(regexp.QuoteMeta(selectRow)).WithArgs(curr.long, curr.short).WillReturnResult(curr.mockResult)
		} else {
			mock.ExpectExec(regexp.QuoteMeta(selectRow)).WithArgs(curr.long, curr.short).WillReturnError(curr.resultErr)
		}

		err = repo.InsertUrl(models.Url{Long: curr.long, Short: curr.short})
		if !errors.Is(err, curr.resultErr) {
			t.Errorf("got different results: %s and %s", err, curr.resultErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
	}
}

func TestGetId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &RepoPostgre{
		db: db,
	}

	testCases := map[string]struct {
		rows      *sqlmock.Rows
		ids       []uint64
		resultErr error
	}{
		"OK": {
			rows:      sqlmock.NewRows([]string{"max"}),
			ids:       []uint64{1},
			resultErr: nil,
		},
		"Bad result": {
			rows:      nil,
			ids:       []uint64{0},
			resultErr: fmt.Errorf("db error"),
		},
	}

	selectRow := "SELECT COUNT(short) + 1 as max FROM url"
	for _, curr := range testCases {
		if curr.resultErr == nil {
			for _, item := range curr.ids {
				curr.rows.AddRow(item)
			}
			mock.ExpectQuery(regexp.QuoteMeta(selectRow)).WithoutArgs().WillReturnRows(curr.rows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(selectRow)).WithoutArgs().WillReturnError(curr.resultErr)
		}

		id, err := repo.GetId()
		if !errors.Is(err, curr.resultErr) {
			t.Errorf("got different results: %s and %s", err, curr.resultErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if id != curr.ids[0] {
			t.Errorf("unexpected result. want %d, got %d", curr.ids[0], id)
			return
		}
	}
}

func TestGetShort(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &RepoPostgre{
		db: db,
	}

	testCases := map[string]struct {
		long       string
		rows       *sqlmock.Rows
		resultRows []string
		resultErr  error
	}{
		"OK": {
			long:       "longg",
			rows:       sqlmock.NewRows([]string{"short"}),
			resultRows: []string{"srt"},
			resultErr:  nil,
		},
		"Bad result": {
			long:       "qwe",
			rows:       nil,
			resultRows: []string{""},
			resultErr:  fmt.Errorf("db error"),
		},
	}

	selectRow := "SELECT short FROM url WHERE long = $1"
	for _, curr := range testCases {
		if curr.resultErr == nil {
			for _, item := range curr.resultRows {
				curr.rows.AddRow(item)
			}
			mock.ExpectQuery(regexp.QuoteMeta(selectRow)).WithArgs(curr.long).WillReturnRows(curr.rows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(selectRow)).WithArgs(curr.long).WillReturnError(curr.resultErr)
		}

		result, err := repo.GetShort(curr.long)
		if !errors.Is(err, curr.resultErr) {
			t.Errorf("got different results: %s and %s", err, curr.resultErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if result != curr.resultRows[0] {
			t.Errorf("unexpected result. want %s, got %s", curr.resultRows[0], result)
			return
		}
	}
}

func TestGetLong(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &RepoPostgre{
		db: db,
	}

	testCases := map[string]struct {
		short      string
		rows       *sqlmock.Rows
		resultRows []string
		resultErr  error
	}{
		"OK": {
			short:      "srt",
			rows:       sqlmock.NewRows([]string{"long"}),
			resultRows: []string{"srt"},
			resultErr:  nil,
		},
		"Bad result": {
			short:      "qwe",
			rows:       nil,
			resultRows: []string{""},
			resultErr:  fmt.Errorf("db error"),
		},
	}

	selectRow := "SELECT long FROM url WHERE short = $1"
	for _, curr := range testCases {
		if curr.resultErr == nil {
			for _, item := range curr.resultRows {
				curr.rows.AddRow(item)
			}
			mock.ExpectQuery(regexp.QuoteMeta(selectRow)).WithArgs(curr.short).WillReturnRows(curr.rows)
		} else {
			mock.ExpectQuery(regexp.QuoteMeta(selectRow)).WithArgs(curr.short).WillReturnError(curr.resultErr)
		}

		result, err := repo.GetLong(curr.short)
		if !errors.Is(err, curr.resultErr) {
			t.Errorf("got different results: %s and %s", err, curr.resultErr)
			return
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
			return
		}
		if result != curr.resultRows[0] {
			t.Errorf("unexpected result. want %s, got %s", curr.resultRows[0], result)
			return
		}
	}
}
