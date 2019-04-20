package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/matrosov-nikita/newsapp/client-service"
)

// ErrInvalidRequestBody happens when request body can not be decoded from JSON.
var ErrInvalidRequestBody = errors.New("could not decode request body")

// ErrInvalidId happens when given id is in wrong format.
var ErrInvalidId = errors.New("invalid id")

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
		h.Error(w, ErrInvalidRequestBody)
		return
	}

	id, err := h.c.CreateNews(form.Header)
	if err != nil {
		h.Error(w, err)
		return
	}

	bs, err := json.Marshal(ResponseNewsForm{
		ID: id,
	})
	if err != nil {
		h.Error(w, fmt.Errorf("fail when marshaling result for news id: %v, get error: %v\n", id, err))
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

	if !h.isValidID(id) {
		h.Error(w, ErrInvalidId)
		return
	}

	news, err := h.c.FindById(id)
	if err != nil {
		if err == client_service.ErrNewsNotFound {
			h.Error(w, client_service.ErrNewsNotFound)
			return
		}

		h.Error(w, err)
		return
	}

	bs, err := json.Marshal(news)
	if err != nil {
		h.Error(w, fmt.Errorf("fail when marshaling result for news id: %v, get error: %v\n", id, err))
		return
	}

	_, err = w.Write(bs)
	if err != nil {
		log.Printf("could not write response for news with id: %v", id)
	}
}

func (h Handler) Error(w http.ResponseWriter, e error) {
	err := customError{Error: e.Error()}
	switch e {
	case ErrInvalidRequestBody, ErrInvalidId:
		err.statusCode = http.StatusBadRequest
	case client_service.ErrNewsNotFound:
		err.statusCode = http.StatusNotFound
	default:
		log.Println(e)
		err.statusCode = http.StatusInternalServerError
		err.Error = "Internal Server Error"
	}

	bs, _ := json.Marshal(err)
	w.WriteHeader(err.statusCode)
	w.Write(bs)
}

type customError struct {
	Error      string `json:"error"`
	statusCode int    `json:"-"`
}

func (h *Handler) isValidID(id string) bool {
	r := regexp.MustCompile("^[0-9a-fA-F]{24}$")
	return r.MatchString(id)
}
