package put

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/model"
	"github.com/nabishec/restapi/internal/storage"
)

type SongPutImp interface {
	PutSongDetail(song *model.Song, songDetail *model.SongDetail) error
}

type Request struct {
	SongData      model.Song       `json:"dataSong" validate:"required"`
	NewSongDetail model.SongDetail `json:"songDetail" validate:"required"`
}

func SongDetail(log *slog.Logger, songPutImp SongPutImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.put.songDetail()"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
			} else {
				log.Error("failed to decode request body", slerr.Err(err))
			}

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.StatusError("bad request"))
			return
		}

		log.Info("request body decoded", slog.Any("songInf: ", req))

		if err := validator.New().Struct(req); err != nil {
			validatorErr := err.(validator.ValidationErrors)
			log.Error("invalid types", slerr.Err(err))

			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, model.StatusError(validatorErr.Error()))
			return
		}

		err = songPutImp.PutSongDetail(&req.SongData, &req.NewSongDetail)
		if errors.Is(err, storage.ErrSongDetailNotFound) {
			log.Info("song detail dosn't exist", slog.String("song:", req.SongData.SongName+
				":"+req.SongData.GroupName))

			w.WriteHeader(http.StatusNotFound) // 404
			render.JSON(w, r, model.StatusError("song detail doesn't exist"))
			return
		}
		if err != nil {
			log.Error("failed add song detail", slerr.Err(err))

			w.WriteHeader(http.StatusInternalServerError) // 500
			render.JSON(w, r, model.StatusError("failed  add song detail"))
			return
		}

		log.Info("song detail changed")
		render.JSON(w, r, model.OK())
	}
}
