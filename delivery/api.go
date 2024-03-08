package delivery

import (
	"log/slog"
	"net/http"
	"test/configs"
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

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		a.lg.Info("get info forbidden method")
		a.sendResponse(w, r, response)
		return
	}

	a.lg.Info("get info done successfully")
	a.sendResponse(w, r, response)
}
