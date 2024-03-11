package delivery

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"test/configs"
	"test/pkg/models"
	"test/pkg/requests"
	"test/usecase"

	"github.com/mailru/easyjson"
)

type API struct {
	core   usecase.ICore
	lg     *slog.Logger
	mx     *http.ServeMux
	adress string
}

func GetApi(c *usecase.Core, l *slog.Logger, cfg *configs.Config) *API {
	l.Info("creating api")
	api := &API{
		core:   c,
		lg:     l.With("module", "api"),
		mx:     http.NewServeMux(),
		adress: cfg.ServerPort,
	}

	api.mx.HandleFunc("/info", api.GetInfo)

	return api
}

func (a *API) ListenAndServe() {
	err := http.ListenAndServe(a.adress, a.mx)
	if err != nil {
		a.lg.Error("listen and serve error", "err", err.Error())
	}
}

func (a *API) sendResponse(w http.ResponseWriter, r *http.Request, response requests.Response) {
	jsonResponse, err := easyjson.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.lg.Error("failed to pack json", "err", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		a.lg.Error("failed to send response", "err", err.Error())
		return
	}
	a.lg.Info("successfully send response")
}

func (a *API) GetInfo(w http.ResponseWriter, r *http.Request) {
	a.lg.Info("new get info request")
	response := requests.Response{Status: http.StatusOK, Body: nil}

	if r.Method == http.MethodGet || r.Method == http.MethodPost {
		var request models.Url
		var responseBody string

		body, err := io.ReadAll(r.Body)
		if err != nil {
			a.lg.Error("get info error", "err", err.Error())
			response.Status = http.StatusBadRequest
			a.sendResponse(w, r, response)
			return
		}
		err = json.Unmarshal(body, &request)
		if err != nil {
			a.lg.Error("get info error", "err", err.Error())
			response.Status = http.StatusBadRequest
			a.sendResponse(w, r, response)
			return
		}

		if r.Method == http.MethodPost {
			responseBody, err = a.core.GetShort(request.Long)
		} else {
			responseBody, err = a.core.GetLong(request.Short)
		}
		if err != nil {
			a.lg.Error("get info error", "err", err.Error())
			switch err {
			case usecase.ErrUncorrectInput:
				response.Status = http.StatusBadRequest
			case usecase.ErrNotFound:
				response.Status = http.StatusNotFound
			default:
				response.Status = http.StatusInternalServerError
			}
			a.sendResponse(w, r, response)
			return
		}

		response.Body = responseBody
		a.lg.Info("get info done successfully")
		a.sendResponse(w, r, response)
		return
	}

	response.Status = http.StatusMethodNotAllowed
	a.lg.Info("get info forbidden method")
	a.sendResponse(w, r, response)
}
