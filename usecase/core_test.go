package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"test/pkg/models"
	"test/repository"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sqids/sqids-go"
)

func TestGetLong(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := repository.NewMockIDbRepo(mockCtrl)
	var buff bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buff, nil))

	core := Core{db: mockObj, lg: logger}

	testCases := map[string]struct {
		short         string
		hasResultErr  bool
		mockResult    string
		mockResultErr error
		times         int
	}{
		"Uncorrect Input": {
			short:         "bad",
			hasResultErr:  true,
			mockResult:    "",
			mockResultErr: nil,
			times:         0,
		},
		"Repo error": {
			short:         "http://somedomen.ru/asdqwezxcf",
			hasResultErr:  true,
			mockResult:    "",
			mockResultErr: errors.New("repo err"),
			times:         1,
		},
		"Not found": {
			short:         "http://somedomen.ru/asdqwezxca",
			hasResultErr:  true,
			mockResult:    "",
			mockResultErr: nil,
			times:         1,
		},
		"OK": {
			short:         "http://somedomen.ru/asdqwezxcz",
			hasResultErr:  false,
			mockResult:    "result",
			mockResultErr: nil,
			times:         1,
		},
	}

	for _, curr := range testCases {
		mockObj.EXPECT().GetLong(curr.short).Return(curr.mockResult, curr.mockResultErr).Times(curr.times)

		res, err := core.GetLong(curr.short)
		if curr.hasResultErr != (err != nil) {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if res != curr.mockResult {
			t.Errorf("results are not the same: want %s, got %s", curr.mockResult, res)
			return
		}
	}
}

func TestGetShort(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := repository.NewMockIDbRepo(mockCtrl)
	var buff bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buff, nil))

	coder, _ := sqids.New(sqids.Options{
		MinLength: 10,
		Alphabet:  AlphabetForShort,
	})

	core := Core{db: mockObj, lg: logger, coder: coder}

	testCases := map[string]struct {
		long              string
		result            string
		hasResultErr      bool
		getShortResult    string
		getShortResultErr error
		getShortTimes     int
		getIdResult       uint64
		getIdResultErr    error
		getIdTimes        int
		insertErr         error
		insertTimes       int
	}{
		"Uncorrect Input": {
			long:         "bad",
			hasResultErr: true,
		},
		"Get short error": {
			long:              "http://something",
			hasResultErr:      true,
			getShortResultErr: fmt.Errorf("repo err"),
			getShortTimes:     1,
		},
		"Get short found": {
			long:              "http://something",
			result:            "found!",
			hasResultErr:      false,
			getShortResult:    "found!",
			getShortResultErr: nil,
			getShortTimes:     1,
		},
		"Get id error": {
			long:              "http://something",
			hasResultErr:      true,
			getShortResultErr: nil,
			getShortTimes:     1,
			getIdResult:       0,
			getIdTimes:        1,
			getIdResultErr:    fmt.Errorf("repo err"),
		},
		"Insert error": {
			long:              "http://something",
			hasResultErr:      true,
			getShortResultErr: nil,
			getShortTimes:     1,
			getIdResult:       1,
			getIdTimes:        1,
			insertErr:         fmt.Errorf("repo err"),
			insertTimes:       1,
		},
		"Ok": {
			long:              "http://something",
			getShortResultErr: nil,
			getShortTimes:     1,
			getIdResult:       1,
			getIdTimes:        1,
			insertErr:         nil,
			insertTimes:       1,
		},
	}

	for _, curr := range testCases {
		mockObj.EXPECT().GetShort(curr.long).Times(curr.getShortTimes).Return(curr.getShortResult, curr.getShortResultErr)
		mockObj.EXPECT().GetId().Times(curr.getIdTimes).Return(curr.getIdResult, curr.getIdResultErr)

		encoded, _ := core.coder.Encode([]uint64{curr.getIdResult})
		if curr.insertTimes > 0 && curr.insertErr == nil {
			curr.result = "http://somedomen.ru/" + encoded
		}
		mockObj.EXPECT().InsertUrl(models.Url{Short: "http://somedomen.ru/" + encoded, Long: curr.long}).Times(curr.insertTimes).Return(curr.insertErr)

		res, err := core.GetShort(curr.long)
		if curr.hasResultErr != (err != nil) {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if res != curr.result {
			t.Errorf("results are different: want %s, got %s", curr.result, res)
			return
		}
	}
}
