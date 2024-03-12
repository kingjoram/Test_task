package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"test/pkg/models"
	"test/pkg/requests"
	"test/usecase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mailru/easyjson"
)

func getResponse(w *httptest.ResponseRecorder) (*requests.Response, error) {
	var response requests.Response

	body, _ := io.ReadAll(w.Body)
	err := easyjson.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("cant unmarshal jsone")
	}

	return &response, nil
}

func createBody(req models.Url) io.Reader {
	jsonReq, _ := json.Marshal(req)

	body := bytes.NewBuffer(jsonReq)
	return body
}

func TestGetInfo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockCore := usecase.NewMockICore(mockCtrl)
	var buff bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buff, nil))

	api := API{core: mockCore, lg: logger}

	testCases := map[string]struct {
		method         string
		body           io.Reader
		result         *requests.Response
		getShortResult string
		getShortErr    error
		getShortTimes  int
		getLongTimes   int
	}{
		"Bad method": {
			method: http.MethodDelete,
			result: &requests.Response{Status: http.StatusMethodNotAllowed, Body: nil},
		},
		"bad request error": {
			method: http.MethodPost,
			result: &requests.Response{Status: http.StatusBadRequest, Body: nil},
			body:   nil,
		},
		"Uncorrect Input error": {
			method:         http.MethodPost,
			result:         &requests.Response{Status: http.StatusBadRequest, Body: nil},
			body:           createBody(models.Url{Long: "input"}),
			getShortResult: "",
			getShortErr:    usecase.ErrUncorrectInput,
			getShortTimes:  1,
		},
		"Not found error": {
			method:         http.MethodGet,
			result:         &requests.Response{Status: http.StatusNotFound, Body: nil},
			body:           createBody(models.Url{Short: "input"}),
			getShortResult: "",
			getShortErr:    usecase.ErrNotFound,
			getLongTimes:   1,
		},
		"Core error": {
			method:         http.MethodPost,
			result:         &requests.Response{Status: http.StatusInternalServerError, Body: nil},
			body:           createBody(models.Url{Long: "input"}),
			getShortResult: "",
			getShortErr:    fmt.Errorf("core err"),
			getShortTimes:  1,
		},
		"Ok": {
			method:         http.MethodPost,
			result:         &requests.Response{Status: http.StatusOK, Body: "result!"},
			body:           createBody(models.Url{Long: "input"}),
			getShortResult: "result!",
			getShortErr:    nil,
			getShortTimes:  1,
		},
	}

	for _, curr := range testCases {
		r := httptest.NewRequest(curr.method, "/info", curr.body)
		w := httptest.NewRecorder()

		mockCore.EXPECT().GetShort("input").Return(curr.getShortResult, curr.getShortErr).Times(curr.getShortTimes)
		mockCore.EXPECT().GetLong("input").Return(curr.getShortResult, curr.getShortErr).Times(curr.getLongTimes)

		api.GetInfo(w, r)

		response, err := getResponse(w)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if response.Status != curr.result.Status {
			t.Errorf("unexpected status: %d, want %d", response.Status, curr.result.Status)
			return
		}
		if !reflect.DeepEqual(response.Body, curr.result.Body) {
			t.Errorf("wanted %v, got %v", curr.result.Body, response.Body)
			return
		}
	}
}
