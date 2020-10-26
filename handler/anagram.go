package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/serebryakov7/aviasales/anagram"
)

type handler struct {
	log    *log.Logger
	mapper *anagram.Anagram
}

func NewHandler(log *log.Logger, m *anagram.Anagram) *handler {
	return &handler{log: log, mapper: m}
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

	res := h.mapper.Find(word)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) {
	words := make([]string, 0)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(data, &words); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.mapper.InsertWords(words...)

	w.WriteHeader(http.StatusNoContent)
}
