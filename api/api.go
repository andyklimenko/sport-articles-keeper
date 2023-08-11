package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/andyklimenko/sport-articles-keeper/config"
	"github.com/andyklimenko/sport-articles-keeper/model"
	"github.com/andyklimenko/sport-articles-keeper/storage"
	"github.com/gorilla/mux"
)

var (
	errParseLimit  = errors.New("can't parse limit")
	errParseOffset = errors.New("can't parse offset")
)

type repo interface {
	GetMany(ctx context.Context, filter storage.ArticleFilter) ([]model.Article, int64, error)
	GetOne(ctx context.Context, id string) (model.Article, error)
}

type API struct {
	repo repo
	srv  *http.Server
}

func (a *API) ListenAndServe() error {
	return a.srv.ListenAndServe()
}

func (api *API) GracefulShutdown(ctx context.Context) {
	slog.Info("graceful shutdown initiated")
	if err := api.srv.Shutdown(ctx); err != nil {
		err = fmt.Errorf("failed to shutdown: %s", err)
		slog.Error("failed to shutdown", slog.Any("error", err))
	}
}

type getOneResponse struct {
	Data model.Article `json:"data"`
}

func (a *API) getOne(w http.ResponseWriter, r *http.Request) {
	id, found := mux.Vars(r)["id"]
	if !found {
		a.respondNotOK(w, http.StatusBadRequest, errors.New("no article id provided"))
		return
	}

	article, err := a.repo.GetOne(r.Context(), id)
	if errors.Is(err, storage.ErrNotFound) {
		a.respondOK(w, http.StatusNotFound, nil)
		return
	}

	if err != nil {
		a.respondNotOK(w, http.StatusBadRequest, fmt.Errorf("looking for article by ID=%s: %w", id, err))
		return
	}

	resp := getOneResponse{
		Data: article,
	}
	a.respondOK(w, http.StatusOK, resp)
}

type getManyResponse struct {
	Data  []model.Article `json:"data"`
	Total int64           `json:"total"`
}

func (a *API) getMany(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, offset, err := parseLimitOffset(queryParams)
	if err != nil {
		a.respondNotOK(w, http.StatusBadRequest, err)
		return
	}

	const defaultLimit = 50
	if limit == 0 {
		limit = defaultLimit
	}

	f := storage.ArticleFilter{
		Limit:  limit,
		Offset: offset,
	}

	articles, total, err := a.repo.GetMany(r.Context(), f)
	if errors.Is(err, storage.ErrNotFound) {
		a.respondOK(w, http.StatusOK, getManyResponse{Data: []model.Article{}})
		return
	}

	if err != nil {
		a.respondNotOK(w, http.StatusInternalServerError, err)
		return
	}

	resp := getManyResponse{
		Data:  articles,
		Total: total,
	}
	a.respondOK(w, http.StatusOK, resp)
}

func setupRouter(a *API) *mux.Router {
	r := mux.NewRouter()

	r.Use(requestLogger, defaultHeaders)

	r.HandleFunc("/articles", a.getMany).Methods(http.MethodGet)
	r.HandleFunc("/articles/{id}", a.getOne).Methods(http.MethodGet)

	return r
}

func New(cfg config.Server, r repo) *API {
	a := API{
		repo: r,
	}

	srv := http.Server{
		Addr:         cfg.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      setupRouter(&a),
	}
	a.srv = &srv

	return &a
}
