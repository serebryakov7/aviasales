package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type finder interface {
	Find(w string) []string
	Insert(words ...string)
}

type handler struct {
	log    *log.Logger
	finder finder
}

func NewHandler(log *log.Logger, m finder) *handler {
	return &handler{log: log, finder: m}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.post(w, r)

	case http.MethodGet:
		h.get(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	if word == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := h.finder.Find(word)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	words := make([]string, 0)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &words); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.finder.Insert(words...)

	w.WriteHeader(http.StatusNoContent)
}
