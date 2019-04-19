package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matrosov-nikita/newsapp/client-service"
)

// ErrInvalidRequestBody happens when request body can not be decoded from JSON.
var ErrInvalidRequestBody = errors.New("could not decode request body")

type Handler struct {
	c      *client_service.NewsClient
	router *mux.Router
}

func NewHandler(c *client_service.NewsClient) *Handler {
	return &Handler{
		c: c,
	}
}

func (h *Handler) Attach(r *mux.Router) {
	r.HandleFunc("/news", h.Create).Methods("POST")
	r.HandleFunc("/news/{id}", h.GetById).Methods("GET")
}

type RequestNewsForm struct {
	Header string `json:"header"`
}

type ResponseNewsForm struct {
	ID string `json:"id"`
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var form RequestNewsForm

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.c.CreateNews(form.Header)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(ResponseNewsForm{
		ID: id,
	})
	if err != nil {
		log.Printf("fail when marshaling result for news id: %v, get error: %v\n", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(bs)
	if err != nil {
		log.Printf("could not write response for news with id: %v\n", id)
	}
}

func (h Handler) GetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	id := mux.Vars(r)["id"]

	news, err := h.c.FindById(id)
	if err != nil {
		if err == client_service.ErrNewsNotFound {
			http.Error(w, client_service.ErrNewsNotFound.Error(), http.StatusNotFound)
			return
		}

		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(news)
	if err != nil {
		log.Printf("fail when marshaling result for news id: %v, get error: %v\n", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bs)
	if err != nil {
		log.Printf("could not write response for news with id: %v", id)
	}
}
