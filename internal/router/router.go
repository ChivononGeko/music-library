package router

import (
	"music-library/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(handler *handlers.SongHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/songs", handler.GetAllSongsHandler).Methods("GET")
	r.HandleFunc("/song/{id}", handler.GetSongHandler).Methods("GET")
	r.HandleFunc("/song", handler.AddSongHandler).Methods("POST")
	r.HandleFunc("/song/{id}", handler.UpdateSongHandler).Methods("PUT")
	r.HandleFunc("/song/{id}", handler.DeleteSongHandler).Methods("DELETE")

	return r
}
